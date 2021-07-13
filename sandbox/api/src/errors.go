package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func ReturnErrors(writer http.ResponseWriter, jsonResponse JSONResponse, message, jsonStatus string, httpStatusCode int) []byte {
	print(message)
	writer.WriteHeader(httpStatusCode)
	jsonResponse.Status = jsonStatus
	jsonResponse.Error = message
	response, err := json.Marshal(&jsonResponse)
	if err != nil {
		writer = ReturnDecodingJSONError(writer, err)
		return response
	}
	return response
}

func ReturnDecodingJSONError(writer http.ResponseWriter, error error) http.ResponseWriter {
	fmt.Println(error.Error())
	http.Error(writer, "Error encoding response object", http.StatusInternalServerError)
	return writer
}

func ShowAndSetError(error error, message string) (status bool) {
	log.Println(message)
	log.Println(error.Error())
	return FAILURE
}
