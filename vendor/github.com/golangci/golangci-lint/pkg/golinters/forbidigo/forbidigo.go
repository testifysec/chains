package forbidigo

import (
	"fmt"
<<<<<<< HEAD
=======
	"sync"
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)

	"github.com/ashanbrown/forbidigo/forbidigo"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
<<<<<<< HEAD
	"github.com/golangci/golangci-lint/pkg/logutils"
=======
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
)

const linterName = "forbidigo"

func New(settings *config.ForbidigoSettings) *goanalysis.Linter {
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
			err := runForbidigo(pass, settings)
=======
			issues, err := runForbidigo(pass, settings)
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

	// Without AnalyzeTypes, LoadModeSyntax is enough.
	// But we cannot make this depend on the settings and have to mirror the mode chosen in GetAllSupportedLinterConfigs,
	// therefore we have to use LoadModeTypesInfo in all cases.
	return goanalysis.NewLinter(
		linterName,
		"Forbids identifiers",
		[]*analysis.Analyzer{analyzer},
		nil,
<<<<<<< HEAD
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func runForbidigo(pass *analysis.Pass, settings *config.ForbidigoSettings) error {
=======
	).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func runForbidigo(pass *analysis.Pass, settings *config.ForbidigoSettings) ([]goanalysis.Issue, error) {
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	options := []forbidigo.Option{
		forbidigo.OptionExcludeGodocExamples(settings.ExcludeGodocExamples),
		// disable "//permit" directives so only "//nolint" directives matters within golangci-lint
		forbidigo.OptionIgnorePermitDirectives(true),
		forbidigo.OptionAnalyzeTypes(settings.AnalyzeTypes),
	}

	// Convert patterns back to strings because that is what NewLinter accepts.
	var patterns []string
	for _, pattern := range settings.Forbid {
		buffer, err := pattern.MarshalString()
		if err != nil {
<<<<<<< HEAD
			return err
		}

=======
			return nil, err
		}
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
		patterns = append(patterns, string(buffer))
	}

	forbid, err := forbidigo.NewLinter(patterns, options...)
	if err != nil {
<<<<<<< HEAD
		return fmt.Errorf("failed to create linter %q: %w", linterName, err)
	}

=======
		return nil, fmt.Errorf("failed to create linter %q: %w", linterName, err)
	}

	var issues []goanalysis.Issue
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	for _, file := range pass.Files {
		runConfig := forbidigo.RunConfig{
			Fset:     pass.Fset,
			DebugLog: logutils.Debug(logutils.DebugKeyForbidigo),
		}
<<<<<<< HEAD

		if settings.AnalyzeTypes {
			runConfig.TypesInfo = pass.TypesInfo
		}

		hints, err := forbid.RunWithConfig(runConfig, file)
		if err != nil {
			return fmt.Errorf("forbidigo linter failed on file %q: %w", file.Name.String(), err)
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
		if settings != nil && settings.AnalyzeTypes {
			runConfig.TypesInfo = pass.TypesInfo
		}
		hints, err := forbid.RunWithConfig(runConfig, file)
		if err != nil {
			return nil, fmt.Errorf("forbidigo linter failed on file %q: %w", file.Name.String(), err)
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
