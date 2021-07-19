package main

import (
	"github.com/davecgh/go-spew/spew"
	"go/parser"
	"go/token"
	"log"
)

type Code struct {
	body   string
	health int
}

func convertToAST(code Code) {
	tokenFileSet := token.NewFileSet()
	ast, error := parser.ParseFile(tokenFileSet, "", code.body, parser.AllErrors)
	if error != nil {
		log.Fatal(error.Error())
	}
	spew.Dump(ast)
}
