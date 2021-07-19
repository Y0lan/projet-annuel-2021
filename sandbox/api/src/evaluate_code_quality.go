package main

import (
	"fmt"
	"go/ast"
	"go/token"
)

func EvaluateCodeQuality(data CodeData) {
	data.Code = minimize(data.Code)
	Ast := convertToAST(data.Code)
	codeComplexity := EvaluateCodeComplexity(Ast)
	fmt.Println(codeComplexity)
}

func EvaluateCodeComplexity(Ast *ast.File) float64 {
	var stats Stats
	fset := token.NewFileSet()
	stats = AnalyzeASTFile(Ast, fset, stats)
	return stats.AverageComplexity()
}
