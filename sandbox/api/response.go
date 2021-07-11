package main

import "net/http"

func setHeader(writer http.ResponseWriter, key, value string) http.ResponseWriter {
	writer.Header().Set(key, value)
	return writer
}
