package main

import (
	"github.com/gorilla/mux"
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

func main() {
	router := mux.NewRouter()
}
