package nolintlint

import (
	"fmt"
<<<<<<< HEAD
=======
	"go/ast"
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	"sync"

	"golang.org/x/tools/go/analysis"

<<<<<<< HEAD
	"github.com/golangci/golangci-lint/pkg/golinters/internal"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	nolintlint "github.com/golangci/golangci-lint/pkg/golinters/nolintlint/internal"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
)

const LinterName = nolintlint.LinterName
=======
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/golinters/nolintlint/internal"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const LinterName = "nolintlint"
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)

func New(settings *config.NoLintLintSettings) *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

<<<<<<< HEAD
	var needs nolintlint.Needs
	if settings.RequireExplanation {
		needs |= nolintlint.NeedsExplanation
	}
	if settings.RequireSpecific {
		needs |= nolintlint.NeedsSpecific
	}
	if !settings.AllowUnused {
		needs |= nolintlint.NeedsUnused
	}

	lnt, err := nolintlint.NewLinter(needs, settings.AllowNoExplanation)
	if err != nil {
		internal.LinterLogger.Fatalf("%s: create analyzer: %v", nolintlint.LinterName, err)
	}

	analyzer := &analysis.Analyzer{
		Name: nolintlint.LinterName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (any, error) {
			issues, err := lnt.Run(pass)
			if err != nil {
				return nil, fmt.Errorf("linter failed to run: %w", err)
=======
	analyzer := &analysis.Analyzer{
		Name: LinterName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (any, error) {
			issues, err := runNoLintLint(pass, settings)
			if err != nil {
				return nil, err
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
			}

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
<<<<<<< HEAD
		nolintlint.LinterName,
=======
		LinterName,
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
		"Reports ill-formed or insufficient nolint directives",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}
<<<<<<< HEAD
=======

func runNoLintLint(pass *analysis.Pass, settings *config.NoLintLintSettings) ([]goanalysis.Issue, error) {
	var needs internal.Needs
	if settings.RequireExplanation {
		needs |= internal.NeedsExplanation
	}
	if settings.RequireSpecific {
		needs |= internal.NeedsSpecific
	}
	if !settings.AllowUnused {
		needs |= internal.NeedsUnused
	}

	lnt, err := internal.NewLinter(needs, settings.AllowNoExplanation)
	if err != nil {
		return nil, err
	}

	nodes := make([]ast.Node, 0, len(pass.Files))
	for _, n := range pass.Files {
		nodes = append(nodes, n)
	}

	lintIssues, err := lnt.Run(pass.Fset, nodes...)
	if err != nil {
		return nil, fmt.Errorf("linter failed to run: %w", err)
	}

	var issues []goanalysis.Issue

	for _, i := range lintIssues {
		expectNoLint := false
		var expectedNolintLinter string
		if ii, ok := i.(internal.UnusedCandidate); ok {
			expectedNolintLinter = ii.ExpectedLinter
			expectNoLint = true
		}

		issue := &result.Issue{
			FromLinter:           LinterName,
			Text:                 i.Details(),
			Pos:                  i.Position(),
			ExpectNoLint:         expectNoLint,
			ExpectedNoLintLinter: expectedNolintLinter,
			Replacement:          i.Replacement(),
		}

		issues = append(issues, goanalysis.NewIssue(issue, pass))
	}

	return issues, nil
}
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
