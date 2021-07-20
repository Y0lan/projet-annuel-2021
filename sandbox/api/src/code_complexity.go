package main

import (
	"fmt"
	"go/ast"
	"go/token"
)

// Stat holds the cyclomatic complexity of a function, along with its package
// and and function name and its position in the source code.
type Stat struct {
	PkgName      string
	FuncName     string
	Complexity   int
	NumberOfLine int
	Pos          token.Position
}

// Stats hold the cyclomatic complexities of many functions.
type Stats []Stat

// AverageComplexity calculates the average cyclomatic complexity of the
// cyclomatic complexities in s.
func (s Stats) AverageComplexity() float64 {
	return float64(s.TotalComplexity()) / float64(len(s))
}

// TotalComplexity calculates the total sum of all cyclomatic
// complexities in s.
func (s Stats) TotalComplexity() uint64 {
	total := uint64(0)
	for _, stat := range s {
		total += uint64(stat.Complexity)
	}
	return total
}

// Complexity calculates the cyclomatic complexity of a function.
// The 'fn' node is either a *ast.FuncDecl or a *ast.FuncLit.
func Complexity(fn ast.Node) int {
	v := complexityVisitor{
		complexity: 1,
	}
	ast.Walk(&v, fn)
	return v.complexity
}

type complexityVisitor struct {
	// complexity is the cyclomatic complexity
	complexity int
}

// Visit implements the ast.Visitor interface.
func (v *complexityVisitor) Visit(n ast.Node) ast.Visitor {
	switch n := n.(type) {
	case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt:
		v.complexity++
	case *ast.CaseClause:
		if n.List != nil { // ignore default case
			v.complexity++
		}
	case *ast.CommClause:
		if n.Comm != nil { // ignore default case
			v.complexity++
		}
	case *ast.BinaryExpr:
		if n.Op == token.LAND || n.Op == token.LOR {
			v.complexity++
		}
	}
	return v
}

// AnalyzeASTFile calculates the cyclomatic complexities of the functions
// and methods in the abstract syntax tree (AST) of a parsed Go file and
// appends the results to the given Stats slice.
func AnalyzeASTFile(f *ast.File, fs *token.FileSet, s Stats) Stats {
	analyzer := &fileAnalyzer{
		file:    f,
		fileSet: fs,
		stats:   s,
	}
	return analyzer.analyze()
}

type fileAnalyzer struct {
	file    *ast.File
	fileSet *token.FileSet
	stats   Stats
}

func (a *fileAnalyzer) analyze() Stats {
	for _, decl := range a.file.Decls {
		a.analyzeDecl(decl)
	}
	return a.stats
}

func (a *fileAnalyzer) analyzeDecl(d ast.Decl) {
	switch decl := d.(type) {
	case *ast.FuncDecl:
		a.addStatIfNotIgnored(decl, funcName(decl))
	case *ast.GenDecl:
		for _, spec := range decl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for _, value := range valueSpec.Values {
				funcLit, ok := value.(*ast.FuncLit)
				if !ok {
					continue
				}
				a.addStatIfNotIgnored(funcLit, valueSpec.Names[0].Name)
			}
		}
	}
}

func (a *fileAnalyzer) addStatIfNotIgnored(node ast.Node, funcName string) {
	a.stats = append(a.stats, Stat{
		PkgName:    a.file.Name.Name,
		FuncName:   funcName,
		Complexity: Complexity(node),
		Pos:        a.fileSet.Position(node.Pos()),
	})
}

// funcName returns the name representation of a function or method:
// "(Type).Name" for methods or simply "Name" for functions.
func funcName(fn *ast.FuncDecl) string {
	if fn.Recv != nil {
		if fn.Recv.NumFields() > 0 {
			typ := fn.Recv.List[0].Type
			return fmt.Sprintf("(%s).%s", recvString(typ), fn.Name)
		}
	}
	return fn.Name.Name
}

// recvString returns a string representation of recv of the
// form "T", "*T", or "BADRECV" (if not a proper receiver type).
func recvString(recv ast.Expr) string {
	switch t := recv.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + recvString(t.X)
	}
	return "BADRECV"
}
