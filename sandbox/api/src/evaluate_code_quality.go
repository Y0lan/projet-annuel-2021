package main

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

func EvaluateCodeQuality(data CodeData, jsonResponse JSONResponse) JSONResponse {
	// MINIMIZE
	data.Code = minimize(data.Code)
	data.Solution = minimize(data.Solution)

	// CONVERT TO AST
	AstSubmitted := convertToAST(data.Code)
	AstSolution := convertToAST(data.Solution)

	// Calculate Code Complexity
	codeComplexitySubmitted := EvaluateCodeComplexity(AstSubmitted)
	codeComplexitySolution := EvaluateCodeComplexity(AstSolution)

	slocSubmitted := GetNumberOfLine(data.Code)
	slocSolution := GetNumberOfLine(data.Solution)

	fmt.Print("cc submit :")
	fmt.Printf("%f", codeComplexitySubmitted)

	fmt.Print("cc solution :")
	fmt.Printf("%f", codeComplexitySolution)
	fmt.Println()

	fmt.Print("lc submit :")
	fmt.Printf("%d", slocSubmitted)
	fmt.Println()

	fmt.Print("lc solution :")
	fmt.Printf("%d", slocSolution)
	fmt.Println()

	jsonResponse.CodeQuality = 0
	return jsonResponse
}

func EvaluateCodeComplexity(Ast *ast.File) float64 {
	var stats Stats
	fset := token.NewFileSet()
	stats = AnalyzeASTFile(Ast, fset, stats)
	return stats.AverageComplexity()
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
