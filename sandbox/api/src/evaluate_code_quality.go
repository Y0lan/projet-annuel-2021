package main

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

func EvaluateCodeQuality(data CodeData, jsonResponse JSONResponse) JSONResponse {
	health := jsonResponse.CodeQuality
	// MINIMIZE
	data.Code = minimize(data.Code)
	data.Solution = minimize(data.Solution)

	// CONVERT TO AST
	astSubmitted := convertToAST(data.Code)
	astSolution := convertToAST(data.Solution)

	health = GetCodeComplexityResult(astSubmitted, astSolution, health)
	health = GetLineNumberResult(data.Code, data.Solution, health)

	jsonResponse.CodeQuality = health

	return jsonResponse
}

func EvaluateCodeComplexity(Ast *ast.File) float64 {
	var stats Stats
	fset := token.NewFileSet()
	stats = AnalyzeASTFile(Ast, fset, stats)
	return stats.AverageComplexity()
}

func GetCodeComplexityResult(astSubmitted, astSolution *ast.File, health float64) float64 {
	// Calculate Code Complexity
	codeComplexitySubmitted := EvaluateCodeComplexity(astSubmitted)
	codeComplexitySolution := EvaluateCodeComplexity(astSolution)
	point := (codeComplexitySubmitted - codeComplexitySolution) * 10
	return health + point
}

// GetNumberOfLine return the number of line excluding function signature and import / package
func GetNumberOfLine(code string) (counter int) {
	scanner := bufio.NewScanner(strings.NewReader(code))
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Trim(line, " \t")
		if !strings.HasPrefix(line, "func") &&
			!strings.HasPrefix(line, "}") &&
			!strings.HasPrefix(line, "\"") &&
			!strings.HasPrefix(line, "import") &&
			!strings.HasPrefix(line, "package") {
			fmt.Println(line)
			counter++
		}
	}
	return
}

func GetLineNumberResult(codeSubmitted, codeSolution string, health float64) float64 {
	slocSubmitted := GetNumberOfLine(codeSubmitted)
	slocSolution := GetNumberOfLine(codeSolution)
	point := float64(slocSolution-slocSubmitted) * 3
	return health + point
}
