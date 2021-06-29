package main

import (
	"github.com/gorilla/mux"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

type language string
type languages []language

type file struct {
	name      string
	extension language
}

func isLangSupported(givenLanguage language) bool {

	supportedLanguages := languages{
		"py",
		"go",
		"rust",
	}

	for _, supportedLanguage := range supportedLanguages {
		if givenLanguage == supportedLanguage {
			return true
		}
	}

	return false
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("file")
	fileName := r.FormValue("file_name")
	if err != nil {
		panic(err)
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			//TODO
		}
	}(file)

	f, err := os.OpenFile("./watch"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0770)
	if err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			//TODO
		}
	}(f)
	_, _ = io.WriteString(w, "File "+fileName+" Uploaded successfully")
	_, _ = io.Copy(f, file)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/compile", UploadFile).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}
