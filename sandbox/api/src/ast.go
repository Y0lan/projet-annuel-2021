package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

func convertToAST(code string) *ast.File {
	tokenFileSet := token.NewFileSet()
	parseFile, error := parser.ParseFile(tokenFileSet, "", code, parser.AllErrors)
	if error != nil {
		log.Fatal(error.Error())
	}
	return parseFile
}
