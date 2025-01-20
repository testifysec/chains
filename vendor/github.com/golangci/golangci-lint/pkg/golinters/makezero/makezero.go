package makezero

import (
	"fmt"
<<<<<<< HEAD
=======
	"sync"
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)

	"github.com/ashanbrown/makezero/makezero"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
<<<<<<< HEAD
=======
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
)

const linterName = "makezero"

func New(settings *config.MakezeroSettings) *goanalysis.Linter {
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
			err := runMakeZero(pass, settings)
=======
			issues, err := runMakeZero(pass, settings)
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
			if err != nil {
				return nil, err
			}

<<<<<<< HEAD
=======
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
		"Finds slice declarations with non-zero initial length",
		[]*analysis.Analyzer{analyzer},
		nil,
<<<<<<< HEAD
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func runMakeZero(pass *analysis.Pass, settings *config.MakezeroSettings) error {
	zero := makezero.NewLinter(settings.Always)

	for _, file := range pass.Files {
		hints, err := zero.Run(pass.Fset, pass.TypesInfo, file)
		if err != nil {
			return fmt.Errorf("makezero linter failed on file %q: %w", file.Name.String(), err)
		}

		for _, hint := range hints {
			pass.Report(analysis.Diagnostic{
				Pos:     hint.Pos(),
				Message: hint.Details(),
			})
		}
	}

	return nil
=======
	).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func runMakeZero(pass *analysis.Pass, settings *config.MakezeroSettings) ([]goanalysis.Issue, error) {
	zero := makezero.NewLinter(settings.Always)

	var issues []goanalysis.Issue

	for _, file := range pass.Files {
		hints, err := zero.Run(pass.Fset, pass.TypesInfo, file)
		if err != nil {
			return nil, fmt.Errorf("makezero linter failed on file %q: %w", file.Name.String(), err)
		}

		for _, hint := range hints {
			issues = append(issues, goanalysis.NewIssue(&result.Issue{
				Pos:        hint.Position(),
				Text:       hint.Details(),
				FromLinter: linterName,
			}, pass))
		}
	}

	return issues, nil
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
}
