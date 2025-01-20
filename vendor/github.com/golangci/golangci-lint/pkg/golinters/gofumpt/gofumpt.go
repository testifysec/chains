package gofumpt

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
<<<<<<< HEAD
=======
	"sync"
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)

	"github.com/shazow/go-diff/difflib"
	"golang.org/x/tools/go/analysis"
	"mvdan.cc/gofumpt/format"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/golinters/internal"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
)

const linterName = "gofumpt"

type differ interface {
	Diff(out io.Writer, a io.ReadSeeker, b io.ReadSeeker) error
}

func New(settings *config.GofumptSettings) *goanalysis.Linter {
<<<<<<< HEAD
=======
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	diff := difflib.New()

	var options format.Options

	if settings != nil {
		options = format.Options{
			LangVersion: getLangVersion(settings),
			ModulePath:  settings.ModulePath,
			ExtraRules:  settings.ExtraRules,
		}
	}

	analyzer := &analysis.Analyzer{
		Name: linterName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run:  goanalysis.DummyRun,
	}

	return goanalysis.NewLinter(
		linterName,
<<<<<<< HEAD
		"Checks if code and import statements are formatted, with additional rules.",
=======
		"Gofumpt checks whether code was gofumpt-ed.",
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		analyzer.Run = func(pass *analysis.Pass) (any, error) {
<<<<<<< HEAD
			err := runGofumpt(lintCtx, pass, diff, options)
=======
			issues, err := runGofumpt(lintCtx, pass, diff, options)
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
			if err != nil {
				return nil, err
			}

<<<<<<< HEAD
			return nil, nil
		}
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}

func runGofumpt(lintCtx *linter.Context, pass *analysis.Pass, diff differ, options format.Options) error {
	for _, file := range pass.Files {
		position, isGoFile := goanalysis.GetGoFilePosition(pass, file)
		if !isGoFile {
			continue
		}

		input, err := os.ReadFile(position.Filename)
		if err != nil {
			return fmt.Errorf("unable to open file %s: %w", position.Filename, err)
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

func runGofumpt(lintCtx *linter.Context, pass *analysis.Pass, diff differ, options format.Options) ([]goanalysis.Issue, error) {
	fileNames := internal.GetFileNames(pass)

	var issues []goanalysis.Issue

	for _, f := range fileNames {
		input, err := os.ReadFile(f)
		if err != nil {
			return nil, fmt.Errorf("unable to open file %s: %w", f, err)
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
		}

		output, err := format.Source(input, options)
		if err != nil {
<<<<<<< HEAD
			return fmt.Errorf("error while running gofumpt: %w", err)
		}

		if !bytes.Equal(input, output) {
			out := bytes.NewBufferString(fmt.Sprintf("--- %[1]s\n+++ %[1]s\n", position.Filename))

			err := diff.Diff(out, bytes.NewReader(input), bytes.NewReader(output))
			if err != nil {
				return fmt.Errorf("error while running gofumpt: %w", err)
			}

			diff := out.String()

			err = internal.ExtractDiagnosticFromPatch(pass, file, diff, lintCtx)
			if err != nil {
				return fmt.Errorf("can't extract issues from gofumpt diff output %q: %w", diff, err)
=======
			return nil, fmt.Errorf("error while running gofumpt: %w", err)
		}

		if !bytes.Equal(input, output) {
			out := bytes.NewBufferString(fmt.Sprintf("--- %[1]s\n+++ %[1]s\n", f))

			err := diff.Diff(out, bytes.NewReader(input), bytes.NewReader(output))
			if err != nil {
				return nil, fmt.Errorf("error while running gofumpt: %w", err)
			}

			diff := out.String()
			is, err := internal.ExtractIssuesFromPatch(diff, lintCtx, linterName, getIssuedTextGoFumpt)
			if err != nil {
				return nil, fmt.Errorf("can't extract issues from gofumpt diff output %q: %w", diff, err)
			}

			for i := range is {
				issues = append(issues, goanalysis.NewIssue(&is[i], pass))
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
			}
		}
	}

<<<<<<< HEAD
	return nil
=======
	return issues, nil
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
}

func getLangVersion(settings *config.GofumptSettings) string {
	if settings == nil || settings.LangVersion == "" {
		// TODO: defaults to "1.15", in the future (v2) must be removed.
		return "go1.15"
	}

	return "go" + strings.TrimPrefix(settings.LangVersion, "go")
}
<<<<<<< HEAD
=======

func getIssuedTextGoFumpt(settings *config.LintersSettings) string {
	text := "File is not `gofumpt`-ed"

	if settings.Gofumpt.ExtraRules {
		text += " with `-extra`"
	}

	return text
}
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
