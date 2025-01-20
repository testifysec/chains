package gochecknoinits

import (
<<<<<<< HEAD
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/golinters/internal"
=======
	"fmt"
	"go/ast"
	"go/token"
	"sync"

	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/golinters/internal"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
)

const linterName = "gochecknoinits"

func New() *goanalysis.Linter {
<<<<<<< HEAD
	analyzer := &analysis.Analyzer{
		Name:     linterName,
		Doc:      goanalysis.TheOnlyanalyzerDoc,
		Run:      run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
=======
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name: linterName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (any, error) {
			var res []goanalysis.Issue
			for _, file := range pass.Files {
				fileIssues := checkFileForInits(file, pass.Fset)
				for i := range fileIssues {
					res = append(res, goanalysis.NewIssue(&fileIssues[i], pass))
				}
			}
			if len(res) == 0 {
				return nil, nil
			}

			mu.Lock()
			resIssues = append(resIssues, res...)
			mu.Unlock()

			return nil, nil
		},
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	}

	return goanalysis.NewLinter(
		linterName,
		"Checks that no init functions are present in Go code",
		[]*analysis.Analyzer{analyzer},
		nil,
<<<<<<< HEAD
	).WithLoadMode(goanalysis.LoadModeSyntax)
}

func run(pass *analysis.Pass) (any, error) {
	insp, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, nil
	}

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(decl ast.Node) {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			return
=======
	).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}

func checkFileForInits(f *ast.File, fset *token.FileSet) []result.Issue {
	var res []result.Issue
	for _, decl := range f.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
		}

		fnName := funcDecl.Name.Name
		if fnName == "init" && funcDecl.Recv.NumFields() == 0 {
<<<<<<< HEAD
			pass.Reportf(funcDecl.Pos(), "don't use %s function", internal.FormatCode(fnName, nil))
		}
	})

	return nil, nil
=======
			res = append(res, result.Issue{
				Pos:        fset.Position(funcDecl.Pos()),
				Text:       fmt.Sprintf("don't use %s function", internal.FormatCode(fnName, nil)),
				FromLinter: linterName,
			})
		}
	}

	return res
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
}
