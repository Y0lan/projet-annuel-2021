package main

import (
	"encoding/json"
	"github.com/CrowdSurge/banner"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const SUCCESS = true
const FAILURE = false

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
	Duration             string `json:"duration"`
	CodeQuality          int    `json:"code-quality"`
}

type CodeResult struct {
	isResultOfTest	 bool
	Duration         string
	CodeQuality      int
	CommandOutput    string
	ExitSuccessfully bool
}

// isLanguageSupported that a language and tells if it's supported by our compiler
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

// createTemporaryFile takes:
// content to write into the file
// extension of the file
// fileName (optional)  so that the file is not created with a random name
// prefix (optional) to put BEFORE the random file name
// suffix (optional) to put AFTER the random file name and BEFORE the extension
func createTemporaryFile(content, extension, fileName, prefix, suffix string) os.File {
	dir := getTemporaryDirectoryForCodeExecution()
	binaryContent := []byte(content)
	var tmpFile *os.File
	var err error

	if fileName == "" {
		tmpFile, err = ioutil.TempFile(dir, prefix+"*"+suffix+"."+extension)
	} else {
		fileName = filepath.Join(dir, fileName+"."+extension)
		tmpFile, err = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
	}
	if err != nil {
		log.Fatal(err)
	}

	defer closeFile(*tmpFile)
	if _, err := tmpFile.Write(binaryContent); err != nil {
		log.Fatal(err)
	}
	return *tmpFile
}

func closeFile(file os.File) {
	if err := file.Close(); err != nil {
		log.Fatal(err)
	}
}

func getTemporaryDirectoryForCodeExecution() string {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		log.Fatal(err.Error())
	}
	return dir
}

func Compile() {
	//TODO
}

func Execute(binaryName string) {
	//TODO
}

func Test() {
	//TODO

}

func ShowAndSetError(error error, message string) (status bool) {
	log.Println(message)
	log.Println(error.Error())
	return false
}

func compilePython(code string) (output string, status bool) {
	pythonFile := createTemporaryFile(code, "py", "", "", "")
	defer os.Remove(pythonFile.Name())

	out, err := exec.Command("python3", pythonFile.Name()).CombinedOutput()
	output = string(out)
	if err != nil {
		status = ShowAndSetError(err, "error in the given python code")
	}
	return
}

func compileRust(code string) (output string, status bool) {

	rustFile := createTemporaryFile(code, "rs", "", "", "")
	defer os.Remove(rustFile.Name())

	out, err := exec.Command("rustc", "--crate-name", "binary", rustFile.Name()).CombinedOutput()
	output = string(out)

	if err != nil {
		status = ShowAndSetError(err, "impossible to compile rust code")
		return
	}

	out, err = exec.Command("./binary").CombinedOutput()
	output = string(out)
	defer os.Remove("binary")

	if err != nil {
		status = ShowAndSetError(err, "impossible to run the binary of the given rust code")
	}

	return
}

func compileGo(code string) (output string, status bool) {

	goFile := createTemporaryFile(code, "go", "", "", "")
	defer os.Remove(goFile.Name())

	out, err := exec.Command("go", "run", goFile.Name()).CombinedOutput()
	output = string(out)

	if err != nil {
		status = ShowAndSetError(err, "impossible to run this go code")
	}
	return
}

func testPython(code, test string) (output string, status bool) {

	pythonCodeFile := createTemporaryFile(code, "py", "main", "", "")
	defer os.Remove(pythonCodeFile.Name())

	pythonTestFile := createTemporaryFile(test, "py", "", "", "")
	defer os.Remove(pythonTestFile.Name())

	out, err := exec.Command("python3", pythonTestFile.Name()).CombinedOutput()
	output = string(out)

	if err != nil {
		status = ShowAndSetError(err, "impossible to run python test")
	}
	return
}

// enableTestFileToAccessSubmittedCode replace EXERCISE_FILE_RANDOM by the name of the file so the test file can access the function
func enableTestFileToAccessSubmittedCode(test, moduleName string) string {
	var replacer = strings.NewReplacer("EXERCISE_FILE_RANDOM", strings.TrimSuffix(filepath.Base(moduleName), path.Ext(moduleName)))
	test = replacer.Replace(test)
	return test
}

func testRust(code, test string) (output string, status bool) {


	rustCodeFile := createTemporaryFile(code, "rs", "", "code", "")
	defer os.Remove(rustCodeFile.Name())

	test = enableTestFileToAccessSubmittedCode(test, rustCodeFile.Name())

	rustTestFile := createTemporaryFile(test, "rs", "", "", "")
	defer os.Remove(rustTestFile.Name())

	out, err := exec.Command("rustc", "--test", "--crate-name", "test", rustTestFile.Name()).CombinedOutput()
	output = string(out)

	if err != nil {
		status = ShowAndSetError(err, "impossible to compile rust test")
		return
	}

	out, err = exec.Command("./test").CombinedOutput()
	output = string(out)
	if err != nil {
		status = ShowAndSetError(err, "failed to run rust test")
	}
	return
}

func testGo(code, test string) (output string, status bool) {

	goModFile := createTemporaryFile("module challenge\n\ngo 1.16\n", "mod", "go", "", "")
	defer os.Remove(goModFile.Name())

	goCodeFile := createTemporaryFile(code, "go", "", "", "")
	defer os.Remove(goCodeFile.Name())

	goTestFile := createTemporaryFile(test, "go", "", "", "_test")
	defer os.Remove(goTestFile.Name())

	out, err := exec.Command("go", "test").CombinedOutput()
	output = string(out)

	if err != nil {
		status = ShowAndSetError(err, "failed to run go test")
	}
	return
}

func populateJsonResponse(result CodeResult, response *JSONResponse) {
	if result.ExitSuccessfully {
		response.Status =
		response.Error =
	}

}

// RedirectToCompiler takes the code in arguments and the language
// and send back the output, and a boolean to know if the program compile successfully
func RedirectToCompiler(code, language string, jsonResponse *JSONResponse) {
	var result CodeResult
	result.isResultOfTest = false
	start := time.Now()
	switch language {
	case "py":
		result.CommandOutput, result.ExitSuccessfully = compilePython(code)
	case "rs":
		result.CommandOutput, result.ExitSuccessfully = compileRust(code)
	case "go":
		result.CommandOutput, result.ExitSuccessfully = compileGo(code)
	}
	result.Duration = time.Since(start).String()
	populateJsonResponse(result, jsonResponse)
	return
}

// RedirectToTester takes the code in arguments and the tests and the language
// and send back the stdout  and a boolean to know if the program tested successfully
func RedirectToTester(code, test, language string, jsonResponse *JSONResponse) (output string, status bool) {
	var result CodeResult
	result.isResultOfTest = true
	switch language {
	case "py":
		result.CommandOutput, result.ExitSuccessfully = testPython(code, test)
	case "rs":
		result.CommandOutput, result.ExitSuccessfully = testRust(code, test)
	case "go":
		result.CommandOutput, result.ExitSuccessfully = testGo(code, test)
	}
	populateJsonResponse(result, jsonResponse)
	return
}

func CompileAndTestCode(writer http.ResponseWriter, request *http.Request) {

	writer = setHeader(writer, "Content-Type", "application/json")
	writer, request, jsonResponse, code := HandleForm(writer, request)

	// Error handling form
	if jsonResponse.Status != "" {
		return
	}

	RedirectToCompiler(code.Code, code.Lang, &jsonResponse)

	// Why launch test on code that doesn't work ?
	if jsonResponse.CompiledSuccessfully {
		RedirectToTester(code.Code, code.Test, code.Lang, &jsonResponse)
	}

	response, err := json.Marshal(&jsonResponse)
	if err != nil {
		writer = ReturnDecodingJSONError(writer, err)
		return
	}
	writer.Write(response)
}

func main() {
	banner.Print("compile box")
	displayIfEveryLangIsCorrectlyConfigured()
	router := mux.NewRouter()
	router.HandleFunc("/compile", CompileAndTestCode)
	log.Fatal(http.ListenAndServe(":8080", router))
}
