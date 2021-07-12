package main

import "net/http"

func HandleForm(writer http.ResponseWriter, request *http.Request) (http.ResponseWriter, *http.Request, JSONResponse, CodeInputByUser) {
	jsonResponse := getDefaultJsonResponse()
	var code CodeInputByUser
	var response []byte

	*request, response = ParseForm(writer, *request, jsonResponse)
	code, response = DecodeForm(writer, *request, jsonResponse)

	if !isLanguageSupported(code.Lang) {
		response = ReturnErrors(writer, jsonResponse, "language "+code.Lang+" is not supported", "error", http.StatusBadRequest)
	}

	writer.Write(response)
	return writer, request, jsonResponse, code
}

func ParseForm(writer http.ResponseWriter, request http.Request, jsonResponse JSONResponse) (http.Request, []byte) {
	err := request.ParseForm()
	if err != nil {
		response := ReturnErrors(writer, jsonResponse, "error parsing input", "error", http.StatusBadRequest)
		return request, response
	}
	return request, nil
}

func DecodeForm(writer http.ResponseWriter, request http.Request, jsonResponse JSONResponse) (CodeInputByUser, []byte) {
	var code CodeInputByUser
	err := decoder.Decode(&code, request.PostForm)

	if err != nil {
		response := ReturnErrors(writer, jsonResponse, "error parsing form", "error", http.StatusBadRequest)
		return code, response
	}
	return code, nil
}
