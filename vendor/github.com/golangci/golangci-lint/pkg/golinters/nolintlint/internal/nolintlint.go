// Package internal provides a linter to ensure that all //nolint directives are followed by explanations
package internal

import (
<<<<<<< HEAD
	"go/token"
	"regexp"
	"strings"

	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/result"
)

const LinterName = "nolintlint"
=======
	"fmt"
	"go/ast"
	"go/token"
	"regexp"
	"strings"
	"unicode"

	"github.com/golangci/golangci-lint/pkg/result"
)

type BaseIssue struct {
	fullDirective                     string
	directiveWithOptionalLeadingSpace string
	position                          token.Position
	replacement                       *result.Replacement
}

//nolint:gocritic // TODO(ldez) must be change in the future.
func (b BaseIssue) Position() token.Position {
	return b.position
}

//nolint:gocritic // TODO(ldez) must be change in the future.
func (b BaseIssue) Replacement() *result.Replacement {
	return b.replacement
}

type ExtraLeadingSpace struct {
	BaseIssue
}

//nolint:gocritic // TODO(ldez) must be change in the future.
func (i ExtraLeadingSpace) Details() string {
	return fmt.Sprintf("directive `%s` should not have more than one leading space", i.fullDirective)
}

func (i ExtraLeadingSpace) String() string { return toString(i) }

type NotMachine struct {
	BaseIssue
}

//nolint:gocritic // TODO(ldez) must be change in the future.
func (i NotMachine) Details() string {
	expected := i.fullDirective[:2] + strings.TrimLeftFunc(i.fullDirective[2:], unicode.IsSpace)
	return fmt.Sprintf("directive `%s` should be written without leading space as `%s`",
		i.fullDirective, expected)
}

func (i NotMachine) String() string { return toString(i) }

type NotSpecific struct {
	BaseIssue
}

//nolint:gocritic // TODO(ldez) must be change in the future.
func (i NotSpecific) Details() string {
	return fmt.Sprintf("directive `%s` should mention specific linter such as `%s:my-linter`",
		i.fullDirective, i.directiveWithOptionalLeadingSpace)
}

func (i NotSpecific) String() string { return toString(i) }

type ParseError struct {
	BaseIssue
}

//nolint:gocritic // TODO(ldez) must be change in the future.
func (i ParseError) Details() string {
	return fmt.Sprintf("directive `%s` should match `%s[:<comma-separated-linters>] [// <explanation>]`",
		i.fullDirective,
		i.directiveWithOptionalLeadingSpace)
}

func (i ParseError) String() string { return toString(i) }

type NoExplanation struct {
	BaseIssue
	fullDirectiveWithoutExplanation string
}

//nolint:gocritic // TODO(ldez) must be change in the future.
func (i NoExplanation) Details() string {
	return fmt.Sprintf("directive `%s` should provide explanation such as `%s // this is why`",
		i.fullDirective, i.fullDirectiveWithoutExplanation)
}

func (i NoExplanation) String() string { return toString(i) }

type UnusedCandidate struct {
	BaseIssue
	ExpectedLinter string
}

//nolint:gocritic // TODO(ldez) must be change in the future.
func (i UnusedCandidate) Details() string {
	details := fmt.Sprintf("directive `%s` is unused", i.fullDirective)
	if i.ExpectedLinter != "" {
		details += fmt.Sprintf(" for linter %q", i.ExpectedLinter)
	}
	return details
}

func (i UnusedCandidate) String() string { return toString(i) }

func toString(issue Issue) string {
	return fmt.Sprintf("%s at %s", issue.Details(), issue.Position())
}

type Issue interface {
	Details() string
	Position() token.Position
	String() string
	Replacement() *result.Replacement
}

type Needs uint
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)

const (
	NeedsMachineOnly Needs = 1 << iota
	NeedsSpecific
	NeedsExplanation
	NeedsUnused
	NeedsAll = NeedsMachineOnly | NeedsSpecific | NeedsExplanation
)

<<<<<<< HEAD
type Needs uint

const commentMark = "//"

=======
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
var commentPattern = regexp.MustCompile(`^//\s*(nolint)(:\s*[\w-]+\s*(?:,\s*[\w-]+\s*)*)?\b`)

// matches a complete nolint directive
var fullDirectivePattern = regexp.MustCompile(`^//\s*nolint(?::(\s*[\w-]+\s*(?:,\s*[\w-]+\s*)*))?\s*(//.*)?\s*\n?$`)

type Linter struct {
	needs           Needs // indicates which linter checks to perform
	excludeByLinter map[string]bool
}

// NewLinter creates a linter that enforces that the provided directives fulfill the provided requirements
func NewLinter(needs Needs, excludes []string) (*Linter, error) {
	excludeByName := make(map[string]bool)
	for _, e := range excludes {
		excludeByName[e] = true
	}

	return &Linter{
		needs:           needs | NeedsMachineOnly,
		excludeByLinter: excludeByName,
	}, nil
}

var (
	leadingSpacePattern      = regexp.MustCompile(`^//(\s*)`)
	trailingBlankExplanation = regexp.MustCompile(`\s*(//\s*)?$`)
)

//nolint:funlen,gocyclo // the function is going to be refactored in the future
<<<<<<< HEAD
func (l Linter) Run(pass *analysis.Pass) ([]goanalysis.Issue, error) {
	var issues []goanalysis.Issue

	for _, file := range pass.Files {
=======
func (l Linter) Run(fset *token.FileSet, nodes ...ast.Node) ([]Issue, error) {
	var issues []Issue

	for _, node := range nodes {
		file, ok := node.(*ast.File)
		if !ok {
			continue
		}

>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
		for _, c := range file.Comments {
			for _, comment := range c.List {
				if !commentPattern.MatchString(comment.Text) {
					continue
				}

				// check for a space between the "//" and the directive
				leadingSpaceMatches := leadingSpacePattern.FindStringSubmatch(comment.Text)

				var leadingSpace string
				if len(leadingSpaceMatches) > 0 {
					leadingSpace = leadingSpaceMatches[1]
				}

<<<<<<< HEAD
				directiveWithOptionalLeadingSpace := commentMark
=======
				directiveWithOptionalLeadingSpace := "//"
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
				if leadingSpace != "" {
					directiveWithOptionalLeadingSpace += " "
				}

<<<<<<< HEAD
				split := strings.Split(strings.SplitN(comment.Text, ":", 2)[0], commentMark)
				directiveWithOptionalLeadingSpace += strings.TrimSpace(split[1])

				pos := pass.Fset.Position(comment.Pos())
				end := pass.Fset.Position(comment.End())

				// check for, report and eliminate leading spaces, so we can check for other issues
				if leadingSpace != "" {
					removeWhitespace := []analysis.SuggestedFix{{
						TextEdits: []analysis.TextEdit{{
							Pos:     token.Pos(pos.Offset),
							End:     token.Pos(pos.Offset + len(commentMark) + len(leadingSpace)),
							NewText: []byte(commentMark),
						}},
					}}

					if (l.needs & NeedsMachineOnly) != 0 {
						issue := &result.Issue{
							FromLinter:     LinterName,
							Text:           formatNotMachine(comment.Text),
							Pos:            pos,
							SuggestedFixes: removeWhitespace,
						}

						issues = append(issues, goanalysis.NewIssue(issue, pass))
					} else if len(leadingSpace) > 1 {
						issue := &result.Issue{
							FromLinter:     LinterName,
							Text:           formatExtraLeadingSpace(comment.Text),
							Pos:            pos,
							SuggestedFixes: removeWhitespace,
						}

						issues = append(issues, goanalysis.NewIssue(issue, pass))
=======
				split := strings.Split(strings.SplitN(comment.Text, ":", 2)[0], "//")
				directiveWithOptionalLeadingSpace += strings.TrimSpace(split[1])

				pos := fset.Position(comment.Pos())
				end := fset.Position(comment.End())

				base := BaseIssue{
					fullDirective:                     comment.Text,
					directiveWithOptionalLeadingSpace: directiveWithOptionalLeadingSpace,
					position:                          pos,
				}

				// check for, report and eliminate leading spaces, so we can check for other issues
				if leadingSpace != "" {
					removeWhitespace := &result.Replacement{
						Inline: &result.InlineFix{
							StartCol:  pos.Column + 1,
							Length:    len(leadingSpace),
							NewString: "",
						},
					}
					if (l.needs & NeedsMachineOnly) != 0 {
						issue := NotMachine{BaseIssue: base}
						issue.BaseIssue.replacement = removeWhitespace
						issues = append(issues, issue)
					} else if len(leadingSpace) > 1 {
						issue := ExtraLeadingSpace{BaseIssue: base}
						issue.BaseIssue.replacement = removeWhitespace
						issue.BaseIssue.replacement.Inline.NewString = " " // assume a single space was intended
						issues = append(issues, issue)
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
					}
				}

				fullMatches := fullDirectivePattern.FindStringSubmatch(comment.Text)
				if len(fullMatches) == 0 {
<<<<<<< HEAD
					issue := &result.Issue{
						FromLinter: LinterName,
						Text:       formatParseError(comment.Text, directiveWithOptionalLeadingSpace),
						Pos:        pos,
					}

					issues = append(issues, goanalysis.NewIssue(issue, pass))

=======
					issues = append(issues, ParseError{BaseIssue: base})
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
					continue
				}

				lintersText, explanation := fullMatches[1], fullMatches[2]

				var linters []string
				if lintersText != "" && !strings.HasPrefix(lintersText, "all") {
					lls := strings.Split(lintersText, ",")
					linters = make([]string, 0, len(lls))
<<<<<<< HEAD
					rangeStart := (pos.Column - 1) + len(commentMark) + len(leadingSpace) + len("nolint:")
=======
					rangeStart := (pos.Column - 1) + len("//") + len(leadingSpace) + len("nolint:")
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
					for i, ll := range lls {
						rangeEnd := rangeStart + len(ll)
						if i < len(lls)-1 {
							rangeEnd++ // include trailing comma
						}
						trimmedLinterName := strings.TrimSpace(ll)
						if trimmedLinterName != "" {
							linters = append(linters, trimmedLinterName)
						}
						rangeStart = rangeEnd
					}
				}

				if (l.needs & NeedsSpecific) != 0 {
					if len(linters) == 0 {
<<<<<<< HEAD
						issue := &result.Issue{
							FromLinter: LinterName,
							Text:       formatNotSpecific(comment.Text, directiveWithOptionalLeadingSpace),
							Pos:        pos,
						}

						issues = append(issues, goanalysis.NewIssue(issue, pass))
=======
						issues = append(issues, NotSpecific{BaseIssue: base})
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
					}
				}

				// when detecting unused directives, we send all the directives through and filter them out in the nolint processor
				if (l.needs & NeedsUnused) != 0 {
<<<<<<< HEAD
					removeNolintCompletely := []analysis.SuggestedFix{{
						TextEdits: []analysis.TextEdit{{
							Pos:     token.Pos(pos.Offset),
							End:     token.Pos(end.Offset),
							NewText: nil,
						}},
					}}

					if len(linters) == 0 {
						issue := &result.Issue{
							FromLinter:     LinterName,
							Text:           formatUnusedCandidate(comment.Text, ""),
							Pos:            pos,
							ExpectNoLint:   true,
							SuggestedFixes: removeNolintCompletely,
						}

						issues = append(issues, goanalysis.NewIssue(issue, pass))
					} else {
						for _, linter := range linters {
							issue := &result.Issue{
								FromLinter:           LinterName,
								Text:                 formatUnusedCandidate(comment.Text, linter),
								Pos:                  pos,
								ExpectNoLint:         true,
								ExpectedNoLintLinter: linter,
							}

							// only offer SuggestedFix if there is a single linter
							// because of issues around commas and the possibility of all
							// linters being removed
							if len(linters) == 1 {
								issue.SuggestedFixes = removeNolintCompletely
							}

							issues = append(issues, goanalysis.NewIssue(issue, pass))
=======
					removeNolintCompletely := &result.Replacement{}

					startCol := pos.Column - 1

					if startCol == 0 {
						// if the directive starts from a new line, remove the line
						removeNolintCompletely.NeedOnlyDelete = true
					} else {
						removeNolintCompletely.Inline = &result.InlineFix{
							StartCol:  startCol,
							Length:    end.Column - pos.Column,
							NewString: "",
						}
					}

					if len(linters) == 0 {
						issue := UnusedCandidate{BaseIssue: base}
						issue.replacement = removeNolintCompletely
						issues = append(issues, issue)
					} else {
						for _, linter := range linters {
							issue := UnusedCandidate{BaseIssue: base, ExpectedLinter: linter}
							// only offer replacement if there is a single linter
							// because of issues around commas and the possibility of all
							// linters being removed
							if len(linters) == 1 {
								issue.replacement = removeNolintCompletely
							}
							issues = append(issues, issue)
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
						}
					}
				}

<<<<<<< HEAD
				if (l.needs&NeedsExplanation) != 0 && (explanation == "" || strings.TrimSpace(explanation) == commentMark) {
=======
				if (l.needs&NeedsExplanation) != 0 && (explanation == "" || strings.TrimSpace(explanation) == "//") {
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
					needsExplanation := len(linters) == 0 // if no linters are mentioned, we must have explanation
					// otherwise, check if we are excluding all the mentioned linters
					for _, ll := range linters {
						if !l.excludeByLinter[ll] { // if a linter does require explanation
							needsExplanation = true
							break
						}
					}

					if needsExplanation {
						fullDirectiveWithoutExplanation := trailingBlankExplanation.ReplaceAllString(comment.Text, "")
<<<<<<< HEAD

						issue := &result.Issue{
							FromLinter: LinterName,
							Text:       formatNoExplanation(comment.Text, fullDirectiveWithoutExplanation),
							Pos:        pos,
						}

						issues = append(issues, goanalysis.NewIssue(issue, pass))
=======
						issues = append(issues, NoExplanation{
							BaseIssue:                       base,
							fullDirectiveWithoutExplanation: fullDirectiveWithoutExplanation,
						})
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
					}
				}
			}
		}
	}

	return issues, nil
}
