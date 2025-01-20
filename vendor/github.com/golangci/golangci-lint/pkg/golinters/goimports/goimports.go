package goimports

import (
	"fmt"
<<<<<<< HEAD
=======
	"sync"
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)

	goimportsAPI "github.com/golangci/gofmt/goimports"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/imports"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/golinters/internal"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
)

const linterName = "goimports"

func New(settings *config.GoImportsSettings) *goanalysis.Linter {
<<<<<<< HEAD
=======
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	analyzer := &analysis.Analyzer{
		Name: linterName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run:  goanalysis.DummyRun,
	}

	return goanalysis.NewLinter(
		linterName,
<<<<<<< HEAD
		"Checks if the code and import statements are formatted according to the 'goimports' command.",
=======
		"Check import statements are formatted according to the 'goimport' command. "+
			"Reformat imports in autofix mode.",
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		imports.LocalPrefix = settings.LocalPrefixes

		analyzer.Run = func(pass *analysis.Pass) (any, error) {
<<<<<<< HEAD
			err := runGoImports(lintCtx, pass)
=======
			issues, err := runGoImports(lintCtx, pass)
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
			if err != nil {
				return nil, err
			}

<<<<<<< HEAD
			return nil, nil
		}
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}

func runGoImports(lintCtx *linter.Context, pass *analysis.Pass) error {
	for _, file := range pass.Files {
		position, isGoFile := goanalysis.GetGoFilePosition(pass, file)
		if !isGoFile {
			continue
		}

		diff, err := goimportsAPI.Run(position.Filename)
		if err != nil { // TODO: skip
			return err
=======
			if len(issues) == 0 {
				return nil, nil
			}

			mu.Lock()
			resIssues = append(resIssues, issues...)
			mu.Unlock()

			return nil, nil
		}
	}).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}

func runGoImports(lintCtx *linter.Context, pass *analysis.Pass) ([]goanalysis.Issue, error) {
	fileNames := internal.GetFileNames(pass)

	var issues []goanalysis.Issue

	for _, f := range fileNames {
		diff, err := goimportsAPI.Run(f)
		if err != nil { // TODO: skip
			return nil, err
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
		}
		if diff == nil {
			continue
		}

<<<<<<< HEAD
		err = internal.ExtractDiagnosticFromPatch(pass, file, string(diff), lintCtx)
		if err != nil {
			return fmt.Errorf("can't extract issues from goimports diff output %q: %w", string(diff), err)
		}
	}

	return nil
=======
		is, err := internal.ExtractIssuesFromPatch(string(diff), lintCtx, linterName, getIssuedTextGoImports)
		if err != nil {
			return nil, fmt.Errorf("can't extract issues from gofmt diff output %q: %w", string(diff), err)
		}

		for i := range is {
			issues = append(issues, goanalysis.NewIssue(&is[i], pass))
		}
	}

	return issues, nil
}

func getIssuedTextGoImports(settings *config.LintersSettings) string {
	text := "File is not `goimports`-ed"

	if settings.Goimports.LocalPrefixes != "" {
		text += " with -local " + settings.Goimports.LocalPrefixes
	}

	return text
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
}
