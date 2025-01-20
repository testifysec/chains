package misspell

import (
	"fmt"
<<<<<<< HEAD
	"go/ast"
	"go/token"
	"strings"
=======
	"go/token"
	"strings"
	"sync"
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	"unicode"

	"github.com/golangci/misspell"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
<<<<<<< HEAD
	"github.com/golangci/golangci-lint/pkg/lint/linter"
=======
	"github.com/golangci/golangci-lint/pkg/golinters/internal"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
)

const linterName = "misspell"

func New(settings *config.MisspellSettings) *goanalysis.Linter {
<<<<<<< HEAD
=======
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	analyzer := &analysis.Analyzer{
		Name: linterName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run:  goanalysis.DummyRun,
	}

	return goanalysis.NewLinter(
		linterName,
		"Finds commonly misspelled English words",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		replacer, ruleErr := createMisspellReplacer(settings)

		analyzer.Run = func(pass *analysis.Pass) (any, error) {
			if ruleErr != nil {
				return nil, ruleErr
			}

<<<<<<< HEAD
			err := runMisspell(lintCtx, pass, replacer, settings.Mode)
=======
			issues, err := runMisspell(lintCtx, pass, replacer, settings.Mode)
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
			if err != nil {
				return nil, err
			}

<<<<<<< HEAD
			return nil, nil
		}
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}

func runMisspell(lintCtx *linter.Context, pass *analysis.Pass, replacer *misspell.Replacer, mode string) error {
	for _, file := range pass.Files {
		err := runMisspellOnFile(lintCtx, pass, file, replacer, mode)
		if err != nil {
			return err
		}
	}

	return nil
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

func runMisspell(lintCtx *linter.Context, pass *analysis.Pass, replacer *misspell.Replacer, mode string) ([]goanalysis.Issue, error) {
	fileNames := internal.GetFileNames(pass)

	var issues []goanalysis.Issue
	for _, filename := range fileNames {
		lintIssues, err := runMisspellOnFile(lintCtx, filename, replacer, mode)
		if err != nil {
			return nil, err
		}

		for i := range lintIssues {
			issues = append(issues, goanalysis.NewIssue(&lintIssues[i], pass))
		}
	}

	return issues, nil
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
}

func createMisspellReplacer(settings *config.MisspellSettings) (*misspell.Replacer, error) {
	replacer := &misspell.Replacer{
		Replacements: misspell.DictMain,
	}

	// Figure out regional variations
	switch strings.ToUpper(settings.Locale) {
	case "":
		// nothing
	case "US":
		replacer.AddRuleList(misspell.DictAmerican)
	case "UK", "GB":
		replacer.AddRuleList(misspell.DictBritish)
	case "NZ", "AU", "CA":
		return nil, fmt.Errorf("unknown locale: %q", settings.Locale)
	}

	err := appendExtraWords(replacer, settings.ExtraWords)
	if err != nil {
		return nil, fmt.Errorf("process extra words: %w", err)
	}

	if len(settings.IgnoreWords) != 0 {
		replacer.RemoveRule(settings.IgnoreWords)
	}

	// It can panic.
	replacer.Compile()

	return replacer, nil
}

<<<<<<< HEAD
func runMisspellOnFile(lintCtx *linter.Context, pass *analysis.Pass, file *ast.File, replacer *misspell.Replacer, mode string) error {
	position, isGoFile := goanalysis.GetGoFilePosition(pass, file)
	if !isGoFile {
		return nil
	}

	fileContent, err := lintCtx.FileCache.GetFileBytes(position.Filename)
	if err != nil {
		return fmt.Errorf("can't get file %s contents: %w", position.Filename, err)
=======
func runMisspellOnFile(lintCtx *linter.Context, filename string, replacer *misspell.Replacer, mode string) ([]result.Issue, error) {
	fileContent, err := lintCtx.FileCache.GetFileBytes(filename)
	if err != nil {
		return nil, fmt.Errorf("can't get file %s contents: %w", filename, err)
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	}

	// `r.ReplaceGo` doesn't find issues inside strings: it searches only inside comments.
	// `r.Replace` searches all words: it treats input as a plain text.
	// The standalone misspell tool uses `r.Replace` by default.
	var replace func(input string) (string, []misspell.Diff)
	switch strings.ToLower(mode) {
	case "restricted":
		replace = replacer.ReplaceGo
	default:
		replace = replacer.Replace
	}

<<<<<<< HEAD
	f := pass.Fset.File(file.Pos())

	_, diffs := replace(string(fileContent))

	for _, diff := range diffs {
		text := fmt.Sprintf("`%s` is a misspelling of `%s`", diff.Original, diff.Corrected)

		start := f.LineStart(diff.Line) + token.Pos(diff.Column)
		end := f.LineStart(diff.Line) + token.Pos(diff.Column+len(diff.Original))

		pass.Report(analysis.Diagnostic{
			Pos:     start,
			End:     end,
			Message: text,
			SuggestedFixes: []analysis.SuggestedFix{{
				TextEdits: []analysis.TextEdit{{
					Pos:     start,
					End:     end,
					NewText: []byte(diff.Corrected),
				}},
			}},
		})
	}

	return nil
=======
	_, diffs := replace(string(fileContent))

	var res []result.Issue

	for _, diff := range diffs {
		text := fmt.Sprintf("`%s` is a misspelling of `%s`", diff.Original, diff.Corrected)

		pos := token.Position{
			Filename: filename,
			Line:     diff.Line,
			Column:   diff.Column + 1,
		}

		replacement := &result.Replacement{
			Inline: &result.InlineFix{
				StartCol:  diff.Column,
				Length:    len(diff.Original),
				NewString: diff.Corrected,
			},
		}

		res = append(res, result.Issue{
			Pos:         pos,
			Text:        text,
			FromLinter:  linterName,
			Replacement: replacement,
		})
	}

	return res, nil
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
}

func appendExtraWords(replacer *misspell.Replacer, extraWords []config.MisspellExtraWords) error {
	if len(extraWords) == 0 {
		return nil
	}

	extra := make([]string, 0, len(extraWords)*2)

	for _, word := range extraWords {
		if word.Typo == "" || word.Correction == "" {
			return fmt.Errorf("typo (%q) and correction (%q) fields should not be empty", word.Typo, word.Correction)
		}

		if strings.ContainsFunc(word.Typo, func(r rune) bool { return !unicode.IsLetter(r) }) {
			return fmt.Errorf("the word %q in the 'typo' field should only contain letters", word.Typo)
		}
		if strings.ContainsFunc(word.Correction, func(r rune) bool { return !unicode.IsLetter(r) }) {
			return fmt.Errorf("the word %q in the 'correction' field should only contain letters", word.Correction)
		}

		extra = append(extra, strings.ToLower(word.Typo), strings.ToLower(word.Correction))
	}

	replacer.AddRuleList(extra)

	return nil
}
