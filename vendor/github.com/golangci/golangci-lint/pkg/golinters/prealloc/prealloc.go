package prealloc

import (
	"fmt"
<<<<<<< HEAD
=======
	"sync"
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)

	"github.com/alexkohler/prealloc/pkg"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/golinters/internal"
<<<<<<< HEAD
=======
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
)

const linterName = "prealloc"

func New(settings *config.PreallocSettings) *goanalysis.Linter {
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
			runPreAlloc(pass, settings)
=======
			issues := runPreAlloc(pass, settings)

			if len(issues) == 0 {
				return nil, nil
			}

			mu.Lock()
			resIssues = append(resIssues, issues...)
			mu.Unlock()
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)

			return nil, nil
		},
	}

	return goanalysis.NewLinter(
		linterName,
		"Finds slice declarations that could potentially be pre-allocated",
		[]*analysis.Analyzer{analyzer},
		nil,
<<<<<<< HEAD
	).WithLoadMode(goanalysis.LoadModeSyntax)
}

func runPreAlloc(pass *analysis.Pass, settings *config.PreallocSettings) {
	hints := pkg.Check(pass.Files, settings.Simple, settings.RangeLoops, settings.ForLoops)

	for _, hint := range hints {
		pass.Report(analysis.Diagnostic{
			Pos:     hint.Pos,
			Message: fmt.Sprintf("Consider pre-allocating %s", internal.FormatCode(hint.DeclaredSliceName, nil)),
		})
	}
=======
	).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}

func runPreAlloc(pass *analysis.Pass, settings *config.PreallocSettings) []goanalysis.Issue {
	var issues []goanalysis.Issue

	hints := pkg.Check(pass.Files, settings.Simple, settings.RangeLoops, settings.ForLoops)

	for _, hint := range hints {
		issues = append(issues, goanalysis.NewIssue(&result.Issue{
			Pos:        pass.Fset.Position(hint.Pos),
			Text:       fmt.Sprintf("Consider pre-allocating %s", internal.FormatCode(hint.DeclaredSliceName, nil)),
			FromLinter: linterName,
		}, pass))
	}

	return issues
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
}
