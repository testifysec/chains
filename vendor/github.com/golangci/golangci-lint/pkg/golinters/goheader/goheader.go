package goheader

import (
	"go/token"
<<<<<<< HEAD
	"strings"
=======
	"sync"
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)

	goheader "github.com/denis-tingaikin/go-header"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
<<<<<<< HEAD
=======
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
)

const linterName = "goheader"

func New(settings *config.GoHeaderSettings) *goanalysis.Linter {
<<<<<<< HEAD
=======
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	conf := &goheader.Configuration{}
	if settings != nil {
		conf = &goheader.Configuration{
			Values:       settings.Values,
			Template:     settings.Template,
			TemplatePath: settings.TemplatePath,
		}
	}

	analyzer := &analysis.Analyzer{
		Name: linterName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (any, error) {
<<<<<<< HEAD
			err := runGoHeader(pass, conf)
=======
			issues, err := runGoHeader(pass, conf)
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
		"Checks if file header matches to pattern",
		[]*analysis.Analyzer{analyzer},
		nil,
<<<<<<< HEAD
	).WithLoadMode(goanalysis.LoadModeSyntax)
}

func runGoHeader(pass *analysis.Pass, conf *goheader.Configuration) error {
	if conf.TemplatePath == "" && conf.Template == "" {
		// User did not pass template, so then do not run go-header linter
		return nil
=======
	).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}

func runGoHeader(pass *analysis.Pass, conf *goheader.Configuration) ([]goanalysis.Issue, error) {
	if conf.TemplatePath == "" && conf.Template == "" {
		// User did not pass template, so then do not run go-header linter
		return nil, nil
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	}

	template, err := conf.GetTemplate()
	if err != nil {
<<<<<<< HEAD
		return err
=======
		return nil, err
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	}

	values, err := conf.GetValues()
	if err != nil {
<<<<<<< HEAD
		return err
=======
		return nil, err
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	}

	a := goheader.New(goheader.WithTemplate(template), goheader.WithValues(values))

<<<<<<< HEAD
	for _, file := range pass.Files {
		position, isGoFile := goanalysis.GetGoFilePosition(pass, file)
		if !isGoFile {
			continue
		}

		issue := a.Analyze(&goheader.Target{File: file, Path: position.Filename})
		if issue == nil {
			continue
		}

		f := pass.Fset.File(file.Pos())

		commentLine := 1
		var offset int

		// Inspired by https://github.com/denis-tingaikin/go-header/blob/4c75a6a2332f025705325d6c71fff4616aedf48f/analyzer.go#L85-L92
		if len(file.Comments) > 0 && file.Comments[0].Pos() < file.Package {
			if !strings.HasPrefix(file.Comments[0].List[0].Text, "/*") {
				// When the comment is "//" there is a one character offset.
				offset = 1
			}
			commentLine = goanalysis.GetFilePositionFor(pass.Fset, file.Comments[0].Pos()).Line
		}

		// Skip issues related to build directives.
		// https://github.com/denis-tingaikin/go-header/issues/18
		if issue.Location().Position-offset < 0 {
			continue
		}

		diag := analysis.Diagnostic{
			Pos:     f.LineStart(issue.Location().Line+1) + token.Pos(issue.Location().Position-offset), // The position of the first divergence.
			Message: issue.Message(),
		}

		if fix := issue.Fix(); fix != nil {
			current := len(fix.Actual)
			for _, s := range fix.Actual {
				current += len(s)
			}

			start := f.LineStart(commentLine)

			end := start + token.Pos(current)

			header := strings.Join(fix.Expected, "\n") + "\n"

			// Adds an extra line between the package and the header.
			if end == file.Package {
				header += "\n"
			}

			diag.SuggestedFixes = []analysis.SuggestedFix{{
				TextEdits: []analysis.TextEdit{{
					Pos:     start,
					End:     end,
					NewText: []byte(header),
				}},
			}}
		}

		pass.Report(diag)
	}

	return nil
=======
	var issues []goanalysis.Issue
	for _, file := range pass.Files {
		path := pass.Fset.Position(file.Pos()).Filename

		i := a.Analyze(&goheader.Target{File: file, Path: path})

		if i == nil {
			continue
		}

		issue := result.Issue{
			Pos: token.Position{
				Line:     i.Location().Line + 1,
				Column:   i.Location().Position,
				Filename: path,
			},
			Text:       i.Message(),
			FromLinter: linterName,
		}

		if fix := i.Fix(); fix != nil {
			issue.LineRange = &result.Range{
				From: issue.Line(),
				To:   issue.Line() + len(fix.Actual) - 1,
			}
			issue.Replacement = &result.Replacement{
				NeedOnlyDelete: len(fix.Expected) == 0,
				NewLines:       fix.Expected,
			}
		}

		issues = append(issues, goanalysis.NewIssue(&issue, pass))
	}

	return issues, nil
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
}
