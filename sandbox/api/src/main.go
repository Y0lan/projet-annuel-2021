package main

import (
	"encoding/json"
	"github.com/CrowdSurge/banner"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"log"
	"net/http"
	"time"
)

const SUCCESS = true
const FAILURE = false

type languages []string

var decoder = schema.NewDecoder()

type JSONResponse struct {
	Status               string  `json:"status"`
	Error                string  `json:"error"`
	Output               string  `json:"output"`
	CompiledSuccessfully bool    `json:"compiled_successfully"`
	TestPassed           bool    `json:"test_passed"`
	Duration             string  `json:"duration"`
	CodeQuality          float64 `json:"code-quality"`
}

type CodeData struct {
	Duration         string
	CodeQuality      float64
	CommandOutput    string
	ExitSuccessfully bool
	Code             string
	Test             string
	Lang             string
	Solution         string
}

type Command struct {
	output string
	status bool
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

// RedirectToCompiler takes the code in arguments and the language
// and send back the output, and a boolean to know if the program compile successfully
func RedirectToCompiler(data CodeData, jsonResponse JSONResponse) JSONResponse {
	start := time.Now()
	command := make(chan Command)
	switch data.Lang {
	case "py":
		go compilePython(data.Code, command)
		result := <-command
		data.CommandOutput, data.ExitSuccessfully = result.output, result.status
	case "rs":
		go compileRust(data.Code, command)
		result := <-command
		data.CommandOutput, data.ExitSuccessfully = result.output, result.status
	case "go":
		data.Code += "\nfunc main() {}"
		go compileGo(data.Code, command)
		result := <-command
		data.CommandOutput, data.ExitSuccessfully = result.output, result.status
	}
	data.Duration = time.Since(start).String()
	jsonResponse = populateJsonResponse(data, jsonResponse)
	return jsonResponse
}

// RedirectToTester takes the code in arguments and the tests and the language
// and send back the stdout  and a boolean to know if the program tested successfully
func RedirectToTester(data CodeData, jsonResponse JSONResponse) JSONResponse {
	command := make(chan Command)
	switch data.Lang {
	case "py":
		go testPython(data.Code, data.Test, command)
		result := <-command
		data.CommandOutput, data.ExitSuccessfully = result.output, result.status
	case "rs":
		go testRust(data.Code, data.Test, command)
		result := <-command
		data.CommandOutput, data.ExitSuccessfully = result.output, result.status
	case "go":
		data.Code += "\nfunc main() {}"
		go testGo(data.Code, data.Test, command)
		result := <-command
		data.CommandOutput, data.ExitSuccessfully = result.output, result.status
	}
	jsonResponse = populateJsonResponse(data, jsonResponse)
	return jsonResponse
}

func CompileAndTestCode(writer http.ResponseWriter, request *http.Request) {

	writer = setHeader(writer, "Content-Type", "application/json")
	writer, request, jsonResponse, data := HandleForm(writer, request)

	// Error handling form
	if jsonResponse.Status != "" {
		return
	}
	jsonResponse = RedirectToCompiler(data, jsonResponse)

	// Why launch test on code that doesn't work ?
	if jsonResponse.CompiledSuccessfully {
		jsonResponse = RedirectToTester(data, jsonResponse)
	}

	// let's analyse the AST
	if jsonResponse.CompiledSuccessfully && jsonResponse.TestPassed {
		jsonResponse = EvaluateCodeQuality(data, jsonResponse)
	}

	response, err := json.Marshal(&jsonResponse)
	if err != nil {
		writer = ReturnDecodingJSONError(writer, err)
		return
	}
	log.Println("200 - SUCCESS")
	_, err = writer.Write(response)
	if err != nil {
		log.Fatal("Could not write response to writer...")
	}
}

func main() {
	banner.Print("compile box")
	displayIfEveryLangIsCorrectlyConfigured()
	router := mux.NewRouter()
	router.HandleFunc("/compile", CompileAndTestCode)
	log.Fatal(http.ListenAndServe(":8080", router))
}
