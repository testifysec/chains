package godox

import (
	"go/token"
	"strings"
<<<<<<< HEAD
=======
	"sync"
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)

	"github.com/matoous/godox"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
<<<<<<< HEAD
=======
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
)

const linterName = "godox"

func New(settings *config.GodoxSettings) *goanalysis.Linter {
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
			runGodox(pass, settings)
=======
			issues := runGodox(pass, settings)

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
		"Tool for detection of FIXME, TODO and other comment keywords",
		[]*analysis.Analyzer{analyzer},
		nil,
<<<<<<< HEAD
	).WithLoadMode(goanalysis.LoadModeSyntax)
}

func runGodox(pass *analysis.Pass, settings *config.GodoxSettings) {
	for _, file := range pass.Files {
		position, isGoFile := goanalysis.GetGoFilePosition(pass, file)
		if !isGoFile {
			continue
		}

		messages := godox.Run(file, pass.Fset, settings.Keywords...)
		if len(messages) == 0 {
			continue
		}

		nonAdjPosition := pass.Fset.PositionFor(file.Pos(), false)

		ft := pass.Fset.File(file.Pos())

		for _, i := range messages {
			pass.Report(analysis.Diagnostic{
				Pos:     ft.LineStart(goanalysis.AdjustPos(i.Pos.Line, nonAdjPosition.Line, position.Line)) + token.Pos(i.Pos.Column),
				Message: strings.TrimRight(i.Message, "\n"),
			})
		}
	}
=======
	).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}

func runGodox(pass *analysis.Pass, settings *config.GodoxSettings) []goanalysis.Issue {
	var messages []godox.Message
	for _, file := range pass.Files {
		messages = append(messages, godox.Run(file, pass.Fset, settings.Keywords...)...)
	}

	if len(messages) == 0 {
		return nil
	}

	issues := make([]goanalysis.Issue, len(messages))

	for k, i := range messages {
		issues[k] = goanalysis.NewIssue(&result.Issue{
			Pos: token.Position{
				Filename: i.Pos.Filename,
				Line:     i.Pos.Line,
			},
			Text:       strings.TrimRight(i.Message, "\n"),
			FromLinter: linterName,
		}, pass)
	}

	return issues
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
}
