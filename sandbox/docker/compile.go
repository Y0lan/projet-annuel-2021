package main

import (
	"errors"
	"github.com/gorilla/mux"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

type languages []string

func isExtensionSupported(givenLanguage string) bool {

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


func getExtensionOf(header multipart.FileHeader) (string, error) {
	var extension string = ""
	isExtension := false
	for char, _ := range header.Filename {
		if isExtension {
			extension += string(rune(char))
		}
		if char == '.' {
			isExtension = true
		}
	}
	if isExtensionSupported(extension) {
		return extension, nil
	}
	return "", errors.New("error: invalid extension")
}

func WriteFile(code multipart.File, header multipart.FileHeader, location string) error  {

	extension, err := getExtensionOf(header)
	if err != nil {
		return err
	}

	defer func(code multipart.File) {
		err := code.Close()
		if err != nil {
			return
		}
	}(code)

	destination, err := os.OpenFile(location + "/unsolved." + extension, os.O_RDONLY|os.O_CREATE, 0770)
	if err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			return
		}
	}(destination)

	_, _ = io.Copy(destination, code)
	return nil
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	//TODO fix this shit
	for file, header, err := range r.FormFile("file") {
		if err != nil {
			panic(err)
		}
		WriteFile(file, *header, "./shared")
	}

	_, _ = io.WriteString(w, "File Uploaded successfully")
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/compile", UploadFile).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}
