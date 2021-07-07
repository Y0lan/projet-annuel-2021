package main

import (
	"encoding/json"
	"fmt"
	"github.com/CrowdSurge/banner"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

type languages []string

var decoder = schema.NewDecoder()

type CodeInputByUser struct {
	Code string
	Test string
	Lang string
}

type JSONResponse struct {
	Status               string `json:"status"`
	Error                string `json:"error"`
	Output               string `json:"output"`
	CompiledSuccessfully bool   `json:"compiled_successfully"`
	TestPassed           bool   `json:"test_passed"`
	Timer                string `json:"timer"`
	CodeQuality          int    `json:"code-quality"`
}

func isLanguageSupported(givenLanguage string) (isSupported bool) {

	isSupported = false
	supportedLanguages := languages{
		"py",
		"go",
		"rs",
	}
	for _, supportedLanguage := range supportedLanguages {
		if givenLanguage == supportedLanguage {
			isSupported = true
			return
		}
	}
	return
}

func createFile(code, extension, fileName string) os.File {
	content := []byte(code)
	var tmpFile *os.File
	var err error
	if fileName == "" {
		tmpFile, err = ioutil.TempFile("", "*."+extension)
	} else {
		dir := os.TempDir()
		fileName = filepath.Join(dir, fileName+"."+extension)
		tmpFile, err = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
	}
	if err != nil {
		log.Fatal(err)
	}

	defer closeFile(*tmpFile)
	if _, err := tmpFile.Write(content); err != nil {
		log.Fatal(err)
	}
	return *tmpFile
}

func closeFile(file os.File) {
	if err := file.Close(); err != nil {
		log.Fatal(err)
	}
}

func compilePython(code string) (output, success string) {
	success = "success"
	pythonFile := createFile(code, "py", "")
	defer os.Remove(pythonFile.Name())
	out, err := exec.Command("python3", pythonFile.Name()).CombinedOutput()
	if err != nil {
		success = "failed"
		log.Println("error in the given python code")
		log.Println(err.Error())
	}
	output = string(out)
	return output, success
}

func compileRust(code string) (output, success string) {
	success = "success"
	rustFile := createFile(code, "rs", "")
	defer os.Remove(rustFile.Name())
	out, err := exec.Command("rustc", "--crate-name", "binary", rustFile.Name()).CombinedOutput()
	if err != nil {
		success = "failed"
		log.Println("impossible to compile rust code")
		log.Println(err.Error())
		return string(out), success
	}

	out, err = exec.Command("./binary").CombinedOutput()
	if err != nil {
		success = "failed"
		log.Println("impossible to run the binary of the given rust code")
		log.Println(err.Error())
		os.Remove("binary")
		output = string(out)
		return
	}

	output = string(out)
	return string(out), success
}

func compileGo(code string) (output, success string) {
	success = "success"
	goFile := createFile(code, "go", "")
	defer os.Remove(goFile.Name())
	out, err := exec.Command("go", "run", goFile.Name()).CombinedOutput()
	if err != nil {
		success = "failed"
		log.Println("impossible to run this go code")
		log.Println(err.Error())
	}
	output = string(out)
	return
}

func testPython(code, test string) (output, success string) {

	success = "success"

	pythonCodeFile := createFile(code, "py", "main")
	defer os.Remove(pythonCodeFile.Name())

	pythonTestFile := createFile(test, "py", "")
	defer os.Remove(pythonTestFile.Name())

	out, err := exec.Command("python3", pythonTestFile.Name()).CombinedOutput()
	if err != nil {
		success = "failed"
		log.Println("impossible to run python test")
		log.Println(err.Error())
	}
	output = string(out)
	return
}

func testRust(code, test string) (output, success string) {
	return
}

func testGo(code, test string) (output, success string) {
	return
}

// compileCode takes the code in arguments and the language
// and send back the output, and a boolean to know if the program compile successfully
func compileCode(code, language string) (output, success string) {
	success = "failed"
	switch language {
	case "py":
		output, success = compilePython(code)
	case "rs":
		output, success = compileRust(code)
	case "go":
		output, success = compileGo(code)
	}
	return
}

// testCode takes the code in arguments and the tests and the language
// and send back the stdout  and a boolean to know if the program tested successfully
func testCode(code, test, language string) (output, success string) {
	success = "failed"
	switch language {
	case "py":
		output, success = testPython(code, test)
	case "rs":
		output, success = testRust(code, test)
	case "go":
		output, success = testGo(code, test)
	}
	return
}

func HandleForm(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	jsonResponse := JSONResponse{
		Status:               "failed",
		CompiledSuccessfully: false,
		TestPassed:           false,
	}
	err := request.ParseForm()
	if err != nil {
		print("Error parsing form!")
		writer.WriteHeader(http.StatusBadRequest)
		jsonResponse.Status = "error"
		response, err := json.Marshal(&jsonResponse)
		if err != nil {
			fmt.Println(err)
			http.Error(writer, "Error encoding response object", http.StatusInternalServerError)
			return
		}
		writer.Write(response)
	}

	var code CodeInputByUser

	err = decoder.Decode(&code, request.PostForm)
	if err != nil {
		print("Error parsing form!")
		writer.WriteHeader(http.StatusBadRequest)
		jsonResponse.Status = "error"
		response, err := json.Marshal(&jsonResponse)
		if err != nil {
			fmt.Println(err)
			http.Error(writer, "Error encoding response object", http.StatusInternalServerError)
			return
		}
		writer.Write(response)
	}

	if !isLanguageSupported(code.Lang) {
		jsonResponse.Status = "error"
		jsonResponse.Error = "language " + code.Lang + " is not supported"
		writer.WriteHeader(http.StatusBadRequest)
		response, err := json.Marshal(&jsonResponse)
		if err != nil {
			fmt.Println(err)
			http.Error(writer, "Error encoding response object", http.StatusInternalServerError)
			return
		}
		writer.Write(response)
		return
	}

	var compilationOutput string
	compilationOutput, jsonResponse.Status = compileCode(code.Code, code.Lang)
	if jsonResponse.Status == "success" {
		jsonResponse.Output = compilationOutput
		jsonResponse.CompiledSuccessfully = true
	}
	if jsonResponse.Status == "failed" {
		jsonResponse.Error = compilationOutput
		jsonResponse.CompiledSuccessfully = false
	}

	var testOutput string
	testOutput, jsonResponse.Status = testCode(code.Code, code.Test, code.Lang)
	if jsonResponse.Status == "success" {
		jsonResponse.Output = compilationOutput + testOutput
		jsonResponse.TestPassed = true
	}
	if jsonResponse.Status == "failed" {
		jsonResponse.Error = compilationOutput + testOutput
		jsonResponse.TestPassed = false
	}

	response, err := json.Marshal(&jsonResponse)
	if err != nil {
		fmt.Println(err)
		http.Error(writer, "Error encoding response object", http.StatusInternalServerError)
		return
	}
	writer.Write(response)
}

func main() {
	banner.Print("compile box")
	out, err := exec.Command("./tests/test_all.sh").Output()
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal("ERROR: DOCKER IMAGE NOT CONFIGURED CORRECTLY. EXITING...")
	}
	fmt.Println(string(out))
	router := mux.NewRouter()
	router.HandleFunc("/compile", HandleForm)
	log.Fatal(http.ListenAndServe(":8080", router))
}