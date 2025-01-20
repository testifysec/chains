package dogsled

import (
<<<<<<< HEAD
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
=======
	"fmt"
	"go/ast"
	"go/token"
	"sync"

	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
)

const linterName = "dogsled"

func New(settings *config.DogsledSettings) *goanalysis.Linter {
<<<<<<< HEAD
=======
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	analyzer := &analysis.Analyzer{
		Name: linterName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (any, error) {
<<<<<<< HEAD
			return run(pass, settings.MaxBlankIdentifiers)
		},
		Requires: []*analysis.Analyzer{inspect.Analyzer},
=======
			issues := runDogsled(pass, settings)

			if len(issues) == 0 {
				return nil, nil
			}

			mu.Lock()
			resIssues = append(resIssues, issues...)
			mu.Unlock()

			return nil, nil
		},
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	}

	return goanalysis.NewLinter(
		linterName,
		"Checks assignments with too many blank identifiers (e.g. x, _, _, _, := f())",
		[]*analysis.Analyzer{analyzer},
		nil,
<<<<<<< HEAD
	).WithLoadMode(goanalysis.LoadModeSyntax)
}

func run(pass *analysis.Pass, maxBlanks int) (any, error) {
	insp, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, nil
	}

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(node ast.Node) {
		funcDecl, ok := node.(*ast.FuncDecl)
		if !ok {
			return
		}

		if funcDecl.Body == nil {
			return
		}

		for _, expr := range funcDecl.Body.List {
			assgnStmt, ok := expr.(*ast.AssignStmt)
			if !ok {
				continue
			}

			numBlank := 0
			for _, left := range assgnStmt.Lhs {
				ident, ok := left.(*ast.Ident)
				if !ok {
					continue
				}
				if ident.Name == "_" {
					numBlank++
				}
			}

			if numBlank > maxBlanks {
				pass.Reportf(assgnStmt.Pos(), "declaration has %v blank identifiers", numBlank)
			}
		}
	})

	return nil, nil
=======
	).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}

func runDogsled(pass *analysis.Pass, settings *config.DogsledSettings) []goanalysis.Issue {
	var reports []goanalysis.Issue
	for _, f := range pass.Files {
		v := &returnsVisitor{
			maxBlanks: settings.MaxBlankIdentifiers,
			f:         pass.Fset,
		}

		ast.Walk(v, f)

		for i := range v.issues {
			reports = append(reports, goanalysis.NewIssue(&v.issues[i], pass))
		}
	}

	return reports
}

type returnsVisitor struct {
	f         *token.FileSet
	maxBlanks int
	issues    []result.Issue
}

func (v *returnsVisitor) Visit(node ast.Node) ast.Visitor {
	funcDecl, ok := node.(*ast.FuncDecl)
	if !ok {
		return v
	}
	if funcDecl.Body == nil {
		return v
	}

	for _, expr := range funcDecl.Body.List {
		assgnStmt, ok := expr.(*ast.AssignStmt)
		if !ok {
			continue
		}

		numBlank := 0
		for _, left := range assgnStmt.Lhs {
			ident, ok := left.(*ast.Ident)
			if !ok {
				continue
			}
			if ident.Name == "_" {
				numBlank++
			}
		}

		if numBlank > v.maxBlanks {
			v.issues = append(v.issues, result.Issue{
				FromLinter: linterName,
				Text:       fmt.Sprintf("declaration has %v blank identifiers", numBlank),
				Pos:        v.f.Position(assgnStmt.Pos()),
			})
		}
	}
	return v
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
}
