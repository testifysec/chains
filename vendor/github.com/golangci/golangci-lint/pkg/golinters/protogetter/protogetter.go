package protogetter

import (
<<<<<<< HEAD
=======
	"sync"

>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	"github.com/ghostiam/protogetter"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
<<<<<<< HEAD
)

func New(settings *config.ProtoGetterSettings) *goanalysis.Linter {
=======
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

func New(settings *config.ProtoGetterSettings) *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	var cfg protogetter.Config
	if settings != nil {
		cfg = protogetter.Config{
			SkipGeneratedBy:         settings.SkipGeneratedBy,
			SkipFiles:               settings.SkipFiles,
			SkipAnyGenerated:        settings.SkipAnyGenerated,
			ReplaceFirstArgInAppend: settings.ReplaceFirstArgInAppend,
		}
	}
<<<<<<< HEAD

	cfg.Mode = protogetter.StandaloneMode

	a := protogetter.NewAnalyzer(&cfg)
=======
	cfg.Mode = protogetter.GolangciLintMode

	a := protogetter.NewAnalyzer(&cfg)
	a.Run = func(pass *analysis.Pass) (any, error) {
		pgIssues, err := protogetter.Run(pass, &cfg)
		if err != nil {
			return nil, err
		}

		issues := make([]goanalysis.Issue, len(pgIssues))
		for i, issue := range pgIssues {
			report := &result.Issue{
				FromLinter: a.Name,
				Pos:        issue.Pos,
				Text:       issue.Message,
				Replacement: &result.Replacement{
					Inline: &result.InlineFix{
						StartCol:  issue.InlineFix.StartCol,
						Length:    issue.InlineFix.Length,
						NewString: issue.InlineFix.NewString,
					},
				},
			}

			issues[i] = goanalysis.NewIssue(report, pass)
		}

		if len(issues) == 0 {
			return nil, nil
		}

		mu.Lock()
		resIssues = append(resIssues, issues...)
		mu.Unlock()

		return nil, nil
	}
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
<<<<<<< HEAD
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
=======
	).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeTypesInfo)
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
}
