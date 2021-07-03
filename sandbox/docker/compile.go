package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"log"
	"net/http"
)

type languages []string

var decoder = schema.NewDecoder()

type CodeInputByUser struct {
	Code string
	Test string
	Lang string
}

type JSONResponse struct {
	Status      string `json:"status"`
	Error       string `json:"error"`
	Stdout      string `json:"stdout"`
	Stderr      string `json:"stderr"`
	TestPassed  bool   `json:"test_passed"`
	Timer       string `json:"timer"`
	CodeQuality int    `json:"code-quality"`
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

func compileCode() (string, string, error) {
	return "", "", nil
}

func runTest() (string, bool, error) {
	return "", false, nil
}

func HandleForm(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	jsonResponse := JSONResponse{
		Status:      "success",
		Error:       "",
		Stdout:      "",
		Stderr:      "",
		TestPassed:  false,
		Timer:       "None", // NO TIME
		CodeQuality: 0,
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

	fmt.Printf("%+v\n", code)
	response, err := json.Marshal(&jsonResponse)

	if err != nil {
		fmt.Println(err)
		http.Error(writer, "Error encoding response object", http.StatusInternalServerError)
		return
	}
	writer.Write(response)
}

func WriteFile(code, extension, location string) error {
	return nil
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/compile", HandleForm)
	log.Fatal(http.ListenAndServe(":8080", router))
}
