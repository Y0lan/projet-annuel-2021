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

func isLanguageSupported(givenLanguage string) bool {

	supportedLanguages := languages{
		"py",
		"go",
		"rs",
	}

	for _, supportedLanguage := range supportedLanguages {
		if givenLanguage == supportedLanguage {
			return true
		}
	}

	return false
}

func createFile(code, extension string) os.File {
	content := []byte(code)
	tmpFile, err := ioutil.TempFile("", "*."+extension)
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

func compilePython(code string) (string, string) {
	success := "success"
	pythonFile := createFile(code, "py")
	defer os.Remove(pythonFile.Name())
	out, err := exec.Command("python3", pythonFile.Name()).CombinedOutput()
	if err != nil {
		success = "failed"
		fmt.Println("error in the given python code")
		fmt.Println(err.Error())
	}
	return string(out), success
}

func compileRust(code string) (string, string) {
	success := "success"
	rustFile := createFile(code, "rs")
	defer os.Remove(rustFile.Name())
	out, err := exec.Command("rustc", "--crate-name", "binary", rustFile.Name()).CombinedOutput()
	if err != nil {
		success = "failed"
		fmt.Println("impossible to compile rust code")
		fmt.Println(err.Error())
		return string(out), success
	}

	out, err = exec.Command("./binary").CombinedOutput()
	if err != nil {
		success = "failed"
		fmt.Println("impossible to run the binary of the given rust code")
		fmt.Println(err.Error())
		os.Remove("binary")
		return string(out), success
	}

	return string(out), success
}

func compileGo(code string) (string, string) {
	success := "success"
	goFile := createFile(code, "go")
	defer os.Remove(goFile.Name())
	out, err := exec.Command("go", "run", goFile.Name()).CombinedOutput()
	if err != nil {
		success = "failed"
		fmt.Println("impossible to run this go code")
		fmt.Println(err.Error())
	}

	return string(out), success
}

func testPython(code string) (string, bool) {
	return "", false
}

func testRust(code string) (string, bool) {
	return "", false
}

func testGo(code string) (string, string, bool) {
	return "", "", false
}

// compileCode takes the code in arguments and the language
// and send back the output, and a boolean to know if the program compile successfully
func compileCode(code, language string) (string, string) {
	output := ""
	success := "failed"
	switch language {
	case "py":
		output, success = compilePython(code)
	case "rs":
		output, success = compileRust(code)
	case "go":
		output, success = compileGo(code)
	}
	return output, success
}

// testCode takes the code in arguments and the tests and the language
// and send back the stdout, stderr, and a boolean to know if the program tested successfully
func testCode(code, test, language string) (string, string, bool) {
	return "", "", false
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

	}

	jsonResponse.Output, jsonResponse.Status = compileCode(code.Code, code.Lang)
	if jsonResponse.Status == "success" {
		jsonResponse.CompiledSuccessfully = true
	} else {
		jsonResponse.CompiledSuccessfully = false
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
