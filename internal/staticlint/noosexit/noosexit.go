package noosexit

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var NoOsExitAnalyzer = &analysis.Analyzer{
	Name: "noosexit",
	Doc:  "check for direct os.Exit calls in main function of main package",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		// Проверяем, является ли пакет "main"
		if pass.Pkg.Name() != "main" {
			continue
		}

		isOsPackExitFunc := func(exprPack string, exprFunc string) bool {
			return exprPack == "os" && exprFunc == "Exit"
		}

		for _, decl := range file.Decls {
			fn, isFn := decl.(*ast.FuncDecl)
			if !isFn || fn.Name.Name != "main" || fn.Recv != nil {
				continue
			}

			// Проверяем тело функции main
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				callExpr, isCall := n.(*ast.CallExpr)
				if !isCall {
					return true
				}

				fun, isIdent := callExpr.Fun.(*ast.SelectorExpr)
				if !isIdent {
					return true
				}

				pkgIdent, isPkg := fun.X.(*ast.Ident)
				if !isPkg {
					return true
				}

				if isOsPackExitFunc(pkgIdent.Name, fun.Sel.Name) {
					pass.Reportf(callExpr.Pos(), "os.Exit call is not allowed in main function")
				}

				return true
			})
		}
	}
	return nil, nil
}
