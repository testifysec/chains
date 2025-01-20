package internal

import (
	"bytes"
	"fmt"
<<<<<<< HEAD
	"go/ast"
	"go/token"
	"slices"
	"strings"

	diffpkg "github.com/sourcegraph/go-diff/diff"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/logutils"
)

type Change struct {
	From, To int
	NewLines []string
=======
	"go/token"
	"strings"

	diffpkg "github.com/sourcegraph/go-diff/diff"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Change struct {
	LineRange   result.Range
	Replacement result.Replacement
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
}

type diffLineType string

const (
	diffLineAdded    diffLineType = "added"
	diffLineOriginal diffLineType = "original"
	diffLineDeleted  diffLineType = "deleted"
)

<<<<<<< HEAD
=======
type fmtTextFormatter func(settings *config.LintersSettings) string

>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
type diffLine struct {
	originalNumber int // 1-based original line number
	typ            diffLineType
	data           string // "+" or "-" stripped line
}

type hunkChangesParser struct {
	// needed because we merge currently added lines with the last original line
	lastOriginalLine *diffLine

	// if the first line of diff is an adding we save all additions to replacementLinesToPrepend
	replacementLinesToPrepend []string

	log logutils.Log

<<<<<<< HEAD
	changes []Change
}

func (p *hunkChangesParser) parse(h *diffpkg.Hunk) []Change {
	lines := parseDiffLines(h)

	for i := 0; i < len(lines); {
		line := lines[i]

		if line.typ == diffLineOriginal {
			p.handleOriginalLine(lines, line, &i)
			continue
		}

		var deletedLines []diffLine
		for ; i < len(lines) && lines[i].typ == diffLineDeleted; i++ {
			deletedLines = append(deletedLines, lines[i])
		}

		var addedLines []string
		for ; i < len(lines) && lines[i].typ == diffLineAdded; i++ {
			addedLines = append(addedLines, lines[i].data)
		}

		if len(deletedLines) != 0 {
			p.handleDeletedLines(deletedLines, addedLines)
			continue
		}

		// no deletions, only additions
		p.handleAddedOnlyLines(addedLines)
	}

	if len(p.replacementLinesToPrepend) != 0 {
		p.log.Infof("The diff contains only additions: no original or deleted lines: %#v", lines)
		return nil
	}

	return p.changes
}

func (p *hunkChangesParser) handleOriginalLine(lines []diffLine, line diffLine, i *int) {
=======
	lines []diffLine

	ret []Change
}

func (p *hunkChangesParser) parseDiffLines(h *diffpkg.Hunk) {
	lines := bytes.Split(h.Body, []byte{'\n'})
	currentOriginalLineNumber := int(h.OrigStartLine)
	var ret []diffLine

	for i, line := range lines {
		dl := diffLine{
			originalNumber: currentOriginalLineNumber,
		}

		lineStr := string(line)

		if strings.HasPrefix(lineStr, "-") {
			dl.typ = diffLineDeleted
			dl.data = strings.TrimPrefix(lineStr, "-")
			currentOriginalLineNumber++
		} else if strings.HasPrefix(lineStr, "+") {
			dl.typ = diffLineAdded
			dl.data = strings.TrimPrefix(lineStr, "+")
		} else {
			if i == len(lines)-1 && lineStr == "" {
				// handle last \n: don't add an empty original line
				break
			}

			dl.typ = diffLineOriginal
			dl.data = strings.TrimPrefix(lineStr, " ")
			currentOriginalLineNumber++
		}

		ret = append(ret, dl)
	}

	// if > 0, then the original file had a 'No newline at end of file' mark
	if h.OrigNoNewlineAt > 0 {
		dl := diffLine{
			originalNumber: currentOriginalLineNumber + 1,
			typ:            diffLineAdded,
			data:           "",
		}
		ret = append(ret, dl)
	}

	p.lines = ret
}

func (p *hunkChangesParser) handleOriginalLine(line diffLine, i *int) {
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	if len(p.replacementLinesToPrepend) == 0 {
		p.lastOriginalLine = &line
		*i++
		return
	}

	// check following added lines for the case:
	// + added line 1
	// original line
	// + added line 2

	*i++
	var followingAddedLines []string
<<<<<<< HEAD
	for ; *i < len(lines) && lines[*i].typ == diffLineAdded; *i++ {
		followingAddedLines = append(followingAddedLines, lines[*i].data)
	}

	change := Change{
		From:     line.originalNumber,
		To:       line.originalNumber,
		NewLines: slices.Concat(p.replacementLinesToPrepend, []string{line.data}, followingAddedLines),
	}
	p.changes = append(p.changes, change)

=======
	for ; *i < len(p.lines) && p.lines[*i].typ == diffLineAdded; *i++ {
		followingAddedLines = append(followingAddedLines, p.lines[*i].data)
	}

	p.ret = append(p.ret, Change{
		LineRange: result.Range{
			From: line.originalNumber,
			To:   line.originalNumber,
		},
		Replacement: result.Replacement{
			NewLines: append(p.replacementLinesToPrepend, append([]string{line.data}, followingAddedLines...)...),
		},
	})
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	p.replacementLinesToPrepend = nil
	p.lastOriginalLine = &line
}

func (p *hunkChangesParser) handleDeletedLines(deletedLines []diffLine, addedLines []string) {
	change := Change{
<<<<<<< HEAD
		From: deletedLines[0].originalNumber,
		To:   deletedLines[len(deletedLines)-1].originalNumber,
	}

	switch {
	case len(addedLines) != 0:
		change.NewLines = slices.Concat(p.replacementLinesToPrepend, addedLines)
		p.replacementLinesToPrepend = nil

	case len(p.replacementLinesToPrepend) != 0:
		// delete-only change with possible prepending
		change.NewLines = slices.Clone(p.replacementLinesToPrepend)
		p.replacementLinesToPrepend = nil
	}

	p.changes = append(p.changes, change)
=======
		LineRange: result.Range{
			From: deletedLines[0].originalNumber,
			To:   deletedLines[len(deletedLines)-1].originalNumber,
		},
	}

	if len(addedLines) != 0 {
		change.Replacement.NewLines = append([]string{}, p.replacementLinesToPrepend...)
		change.Replacement.NewLines = append(change.Replacement.NewLines, addedLines...)
		if len(p.replacementLinesToPrepend) != 0 {
			p.replacementLinesToPrepend = nil
		}

		p.ret = append(p.ret, change)
		return
	}

	// delete-only change with possible prepending
	if len(p.replacementLinesToPrepend) != 0 {
		change.Replacement.NewLines = p.replacementLinesToPrepend
		p.replacementLinesToPrepend = nil
	} else {
		change.Replacement.NeedOnlyDelete = true
	}

	p.ret = append(p.ret, change)
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
}

func (p *hunkChangesParser) handleAddedOnlyLines(addedLines []string) {
	if p.lastOriginalLine == nil {
		// the first line is added; the diff looks like:
		// 1. + ...
		// 2. - ...
		// or
		// 1. + ...
		// 2. ...

		p.replacementLinesToPrepend = addedLines
<<<<<<< HEAD

=======
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
		return
	}

	// add-only change merged into the last original line with possible prepending
<<<<<<< HEAD
	change := Change{
		From:     p.lastOriginalLine.originalNumber,
		To:       p.lastOriginalLine.originalNumber,
		NewLines: slices.Concat(p.replacementLinesToPrepend, []string{p.lastOriginalLine.data}, addedLines),
	}

	p.changes = append(p.changes, change)

	p.replacementLinesToPrepend = nil
}

func parseDiffLines(h *diffpkg.Hunk) []diffLine {
	lines := bytes.Split(h.Body, []byte{'\n'})

	currentOriginalLineNumber := int(h.OrigStartLine)

	var diffLines []diffLine

	for i, line := range lines {
		dl := diffLine{
			originalNumber: currentOriginalLineNumber,
		}

		if i == len(lines)-1 && len(line) == 0 {
			// handle last \n: don't add an empty original line
			break
		}

		lineStr := string(line)

		switch {
		case strings.HasPrefix(lineStr, "-"):
			dl.typ = diffLineDeleted
			dl.data = strings.TrimPrefix(lineStr, "-")
			currentOriginalLineNumber++

		case strings.HasPrefix(lineStr, "+"):
			dl.typ = diffLineAdded
			dl.data = strings.TrimPrefix(lineStr, "+")

		default:
			dl.typ = diffLineOriginal
			dl.data = strings.TrimPrefix(lineStr, " ")
			currentOriginalLineNumber++
		}

		diffLines = append(diffLines, dl)
	}

	// if > 0, then the original file had a 'No newline at end of file' mark
	if h.OrigNoNewlineAt > 0 {
		dl := diffLine{
			originalNumber: currentOriginalLineNumber + 1,
			typ:            diffLineAdded,
			data:           "",
		}
		diffLines = append(diffLines, dl)
	}

	return diffLines
}

func ExtractDiagnosticFromPatch(
	pass *analysis.Pass,
	file *ast.File,
	patch string,
	lintCtx *linter.Context,
) error {
	diffs, err := diffpkg.ParseMultiFileDiff([]byte(patch))
	if err != nil {
		return fmt.Errorf("can't parse patch: %w", err)
	}

	if len(diffs) == 0 {
		return fmt.Errorf("got no diffs from patch parser: %v", patch)
	}

	ft := pass.Fset.File(file.Pos())

	adjLine := pass.Fset.PositionFor(file.Pos(), false).Line - pass.Fset.PositionFor(file.Pos(), true).Line

=======
	p.ret = append(p.ret, Change{
		LineRange: result.Range{
			From: p.lastOriginalLine.originalNumber,
			To:   p.lastOriginalLine.originalNumber,
		},
		Replacement: result.Replacement{
			NewLines: append(p.replacementLinesToPrepend, append([]string{p.lastOriginalLine.data}, addedLines...)...),
		},
	})
	p.replacementLinesToPrepend = nil
}

func (p *hunkChangesParser) parse(h *diffpkg.Hunk) []Change {
	p.parseDiffLines(h)

	for i := 0; i < len(p.lines); {
		line := p.lines[i]
		if line.typ == diffLineOriginal {
			p.handleOriginalLine(line, &i)
			continue
		}

		var deletedLines []diffLine
		for ; i < len(p.lines) && p.lines[i].typ == diffLineDeleted; i++ {
			deletedLines = append(deletedLines, p.lines[i])
		}

		var addedLines []string
		for ; i < len(p.lines) && p.lines[i].typ == diffLineAdded; i++ {
			addedLines = append(addedLines, p.lines[i].data)
		}

		if len(deletedLines) != 0 {
			p.handleDeletedLines(deletedLines, addedLines)
			continue
		}

		// no deletions, only additions
		p.handleAddedOnlyLines(addedLines)
	}

	if len(p.replacementLinesToPrepend) != 0 {
		p.log.Infof("The diff contains only additions: no original or deleted lines: %#v", p.lines)
		return nil
	}

	return p.ret
}

func ExtractIssuesFromPatch(patch string, lintCtx *linter.Context, linterName string, formatter fmtTextFormatter) ([]result.Issue, error) {
	diffs, err := diffpkg.ParseMultiFileDiff([]byte(patch))
	if err != nil {
		return nil, fmt.Errorf("can't parse patch: %w", err)
	}

	if len(diffs) == 0 {
		return nil, fmt.Errorf("got no diffs from patch parser: %v", patch)
	}

	var issues []result.Issue
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	for _, d := range diffs {
		if len(d.Hunks) == 0 {
			lintCtx.Log.Warnf("Got no hunks in diff %+v", d)
			continue
		}

		for _, hunk := range d.Hunks {
			p := hunkChangesParser{log: lintCtx.Log}

			changes := p.parse(hunk)

			for _, change := range changes {
<<<<<<< HEAD
				pass.Report(toDiagnostic(ft, change, adjLine))
=======
				i := result.Issue{
					FromLinter: linterName,
					Pos: token.Position{
						Filename: d.NewName,
						Line:     change.LineRange.From,
					},
					Text:        formatter(lintCtx.Settings()),
					Replacement: &change.Replacement,
				}
				if change.LineRange.From != change.LineRange.To {
					i.LineRange = &change.LineRange
				}

				issues = append(issues, i)
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
			}
		}
	}

<<<<<<< HEAD
	return nil
}

func toDiagnostic(ft *token.File, change Change, adjLine int) analysis.Diagnostic {
	from := change.From + adjLine
	if from > ft.LineCount() {
		from = ft.LineCount()
	}

	start := ft.LineStart(from)

	end := goanalysis.EndOfLinePos(ft, change.To+adjLine)

	return analysis.Diagnostic{
		Pos:     start,
		End:     end,
		Message: "File is not properly formatted",
		SuggestedFixes: []analysis.SuggestedFix{{
			TextEdits: []analysis.TextEdit{{
				Pos:     start,
				End:     end,
				NewText: []byte(strings.Join(change.NewLines, "\n")),
			}},
		}},
	}
=======
	return issues, nil
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
}
