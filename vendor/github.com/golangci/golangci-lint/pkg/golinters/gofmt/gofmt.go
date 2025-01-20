package gofmt

import (
	"fmt"
<<<<<<< HEAD
=======
	"sync"
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)

	gofmtAPI "github.com/golangci/gofmt/gofmt"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/golinters/internal"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
)

const linterName = "gofmt"

func New(settings *config.GoFmtSettings) *goanalysis.Linter {
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
		"Checks if the code is formatted according to 'gofmt' command.",
=======
		"Gofmt checks whether code was gofmt-ed. By default "+
			"this tool runs with -s option to check for code simplification",
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		analyzer.Run = func(pass *analysis.Pass) (any, error) {
<<<<<<< HEAD
			err := runGofmt(lintCtx, pass, settings)
=======
			issues, err := runGofmt(lintCtx, pass, settings)
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
			if err != nil {
				return nil, err
			}

<<<<<<< HEAD
			return nil, nil
		}
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}

func runGofmt(lintCtx *linter.Context, pass *analysis.Pass, settings *config.GoFmtSettings) error {
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

func runGofmt(lintCtx *linter.Context, pass *analysis.Pass, settings *config.GoFmtSettings) ([]goanalysis.Issue, error) {
	fileNames := internal.GetFileNames(pass)

>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	var rewriteRules []gofmtAPI.RewriteRule
	for _, rule := range settings.RewriteRules {
		rewriteRules = append(rewriteRules, gofmtAPI.RewriteRule(rule))
	}

<<<<<<< HEAD
	for _, file := range pass.Files {
		position, isGoFile := goanalysis.GetGoFilePosition(pass, file)
		if !isGoFile {
			continue
		}

		diff, err := gofmtAPI.RunRewrite(position.Filename, settings.Simplify, rewriteRules)
		if err != nil { // TODO: skip
			return err
=======
	var issues []goanalysis.Issue

	for _, f := range fileNames {
		diff, err := gofmtAPI.RunRewrite(f, settings.Simplify, rewriteRules)
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
			return fmt.Errorf("can't extract issues from gofmt diff output %q: %w", string(diff), err)
		}
	}

	return nil
=======
		is, err := internal.ExtractIssuesFromPatch(string(diff), lintCtx, linterName, getIssuedTextGoFmt)
		if err != nil {
			return nil, fmt.Errorf("can't extract issues from gofmt diff output %q: %w", string(diff), err)
		}

		for i := range is {
			issues = append(issues, goanalysis.NewIssue(&is[i], pass))
		}
	}

	return issues, nil
}

func getIssuedTextGoFmt(settings *config.LintersSettings) string {
	text := "File is not `gofmt`-ed"
	if settings.Gofmt.Simplify {
		text += " with `-s`"
	}
	for _, rule := range settings.Gofmt.RewriteRules {
		text += fmt.Sprintf(" `-r '%s -> %s'`", rule.Pattern, rule.Replacement)
	}

	return text
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
}
