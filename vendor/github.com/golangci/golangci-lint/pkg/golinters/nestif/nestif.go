package nestif

import (
<<<<<<< HEAD
=======
	"sort"
	"sync"

>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	"github.com/nakabonne/nestif"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
<<<<<<< HEAD
=======
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
)

const linterName = "nestif"

func New(settings *config.NestifSettings) *goanalysis.Linter {
<<<<<<< HEAD
	analyzer := &analysis.Analyzer{
		Name: linterName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (any, error) {
			runNestIf(pass, settings)
=======
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name: goanalysis.TheOnlyAnalyzerName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (any, error) {
			issues := runNestIf(pass, settings)

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
		"Reports deeply nested if statements",
		[]*analysis.Analyzer{analyzer},
		nil,
<<<<<<< HEAD
	).WithLoadMode(goanalysis.LoadModeSyntax)
}

func runNestIf(pass *analysis.Pass, settings *config.NestifSettings) {
=======
	).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}

func runNestIf(pass *analysis.Pass, settings *config.NestifSettings) []goanalysis.Issue {
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	checker := &nestif.Checker{
		MinComplexity: settings.MinComplexity,
	}

<<<<<<< HEAD
	for _, file := range pass.Files {
		position, isGoFile := goanalysis.GetGoFilePosition(pass, file)
		if !isGoFile {
			continue
		}

		issues := checker.Check(file, pass.Fset)
		if len(issues) == 0 {
			continue
		}

		nonAdjPosition := pass.Fset.PositionFor(file.Pos(), false)

		f := pass.Fset.File(file.Pos())

		for _, issue := range issues {
			pass.Report(analysis.Diagnostic{
				Pos:     f.LineStart(goanalysis.AdjustPos(issue.Pos.Line, nonAdjPosition.Line, position.Line)),
				Message: issue.Message,
			})
		}
	}
=======
	var lintIssues []nestif.Issue
	for _, f := range pass.Files {
		lintIssues = append(lintIssues, checker.Check(f, pass.Fset)...)
	}

	if len(lintIssues) == 0 {
		return nil
	}

	sort.SliceStable(lintIssues, func(i, j int) bool {
		return lintIssues[i].Complexity > lintIssues[j].Complexity
	})

	issues := make([]goanalysis.Issue, 0, len(lintIssues))
	for _, i := range lintIssues {
		issues = append(issues, goanalysis.NewIssue(&result.Issue{
			Pos:        i.Pos,
			Text:       i.Message,
			FromLinter: linterName,
		}, pass))
	}

	return issues
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
}
