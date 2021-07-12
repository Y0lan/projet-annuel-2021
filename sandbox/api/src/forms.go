package main

import (
	"log"
	"net/http"
)

func HandleForm(writer http.ResponseWriter, request *http.Request) (http.ResponseWriter, *http.Request, JSONResponse, CodeData) {
	jsonResponse := getDefaultJsonResponse()
	var data CodeData
	var response []byte

	*request, response = ParseForm(writer, *request, jsonResponse)
	data, response = DecodeForm(writer, *request, jsonResponse)

	if !isLanguageSupported(data.Lang) {
		response = ReturnErrors(writer, jsonResponse, "language "+data.Lang+" is not supported", "error", http.StatusBadRequest)
	}

	_, err := writer.Write(response)
	if err != nil {
		log.Fatal("Error while handling form, wtf did you do ???")
	}
	return writer, request, jsonResponse, data
}

func ParseForm(writer http.ResponseWriter, request http.Request, jsonResponse JSONResponse) (http.Request, []byte) {
	err := request.ParseForm()
	if err != nil {
		response := ReturnErrors(writer, jsonResponse, "error parsing input", "error", http.StatusBadRequest)
		return request, response
	}
	return request, nil
}

func DecodeForm(writer http.ResponseWriter, request http.Request, jsonResponse JSONResponse) (CodeData, []byte) {
	var data CodeData
	err := decoder.Decode(&data, request.PostForm)

	if err != nil {
		response := ReturnErrors(writer, jsonResponse, "error parsing form", "error", http.StatusBadRequest)
		return data, response
	}
	return data, nil
}
