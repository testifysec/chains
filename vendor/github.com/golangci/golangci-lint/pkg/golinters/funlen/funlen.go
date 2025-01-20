package funlen

import (
<<<<<<< HEAD
=======
	"go/token"
	"strings"
	"sync"

>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	"github.com/ultraware/funlen"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
<<<<<<< HEAD
)

type Config struct {
	lineLimit      int
	stmtLimit      int
	ignoreComments bool
}

func New(settings *config.FunlenSettings) *goanalysis.Linter {
	cfg := Config{}
	if settings != nil {
		cfg.lineLimit = settings.Lines
		cfg.stmtLimit = settings.Statements
		cfg.ignoreComments = !settings.IgnoreComments
	}

	a := funlen.NewAnalyzer(cfg.lineLimit, cfg.stmtLimit, cfg.ignoreComments)

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
=======
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const linterName = "funlen"

func New(settings *config.FunlenSettings) *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name: linterName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (any, error) {
			issues := runFunlen(pass, settings)

			if len(issues) == 0 {
				return nil, nil
			}

			mu.Lock()
			resIssues = append(resIssues, issues...)
			mu.Unlock()

			return nil, nil
		},
	}

	return goanalysis.NewLinter(
		linterName,
		"Tool for detection of long functions",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}

func runFunlen(pass *analysis.Pass, settings *config.FunlenSettings) []goanalysis.Issue {
	var lintIssues []funlen.Message
	for _, file := range pass.Files {
		fileIssues := funlen.Run(file, pass.Fset, settings.Lines, settings.Statements, settings.IgnoreComments)
		lintIssues = append(lintIssues, fileIssues...)
	}

	if len(lintIssues) == 0 {
		return nil
	}

	issues := make([]goanalysis.Issue, len(lintIssues))
	for k, i := range lintIssues {
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
