package godot

import (
<<<<<<< HEAD
	"cmp"
=======
	"sync"
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)

	"github.com/tetafro/godot"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
<<<<<<< HEAD
=======
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
)

const linterName = "godot"

func New(settings *config.GodotSettings) *goanalysis.Linter {
<<<<<<< HEAD
=======
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	var dotSettings godot.Settings

	if settings != nil {
		dotSettings = godot.Settings{
			Scope:   godot.Scope(settings.Scope),
			Exclude: settings.Exclude,
			Period:  settings.Period,
			Capital: settings.Capital,
		}

		// Convert deprecated setting
<<<<<<< HEAD
		if settings.CheckAll != nil && *settings.CheckAll {
=======
		if settings.CheckAll {
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
			dotSettings.Scope = godot.AllScope
		}
	}

<<<<<<< HEAD
	dotSettings.Scope = cmp.Or(dotSettings.Scope, godot.DeclScope)
=======
	if dotSettings.Scope == "" {
		dotSettings.Scope = godot.DeclScope
	}
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)

	analyzer := &analysis.Analyzer{
		Name: linterName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (any, error) {
<<<<<<< HEAD
			err := runGodot(pass, dotSettings)
=======
			issues, err := runGodot(pass, dotSettings)
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
		"Check if comments end in a period",
		[]*analysis.Analyzer{analyzer},
		nil,
<<<<<<< HEAD
	).WithLoadMode(goanalysis.LoadModeSyntax)
}

func runGodot(pass *analysis.Pass, settings godot.Settings) error {
	for _, file := range pass.Files {
		iss, err := godot.Run(file, pass.Fset, settings)
		if err != nil {
			return err
		}

		if len(iss) == 0 {
			continue
		}

		f := pass.Fset.File(file.Pos())

		for _, i := range iss {
			start := f.Pos(i.Pos.Offset)
			end := goanalysis.EndOfLinePos(f, i.Pos.Line)

			pass.Report(analysis.Diagnostic{
				Pos:     start,
				End:     end,
				Message: i.Message,
				SuggestedFixes: []analysis.SuggestedFix{{
					TextEdits: []analysis.TextEdit{{
						Pos:     start,
						End:     end,
						NewText: []byte(i.Replacement),
					}},
				}},
			})
		}
	}

	return nil
=======
	).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}

func runGodot(pass *analysis.Pass, settings godot.Settings) ([]goanalysis.Issue, error) {
	var lintIssues []godot.Issue
	for _, file := range pass.Files {
		iss, err := godot.Run(file, pass.Fset, settings)
		if err != nil {
			return nil, err
		}
		lintIssues = append(lintIssues, iss...)
	}

	if len(lintIssues) == 0 {
		return nil, nil
	}

	issues := make([]goanalysis.Issue, len(lintIssues))
	for k, i := range lintIssues {
		issue := result.Issue{
			Pos:        i.Pos,
			Text:       i.Message,
			FromLinter: linterName,
			Replacement: &result.Replacement{
				NewLines: []string{i.Replacement},
			},
		}

		issues[k] = goanalysis.NewIssue(&issue, pass)
	}

	return issues, nil
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
}
