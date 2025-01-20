<<<<<<< HEAD
// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// This file is inspired by go/analysis/internal/checker/checker.go

package processors

import (
	"errors"
	"fmt"
	"os"
	"slices"

	"golang.org/x/exp/maps"

	"github.com/golangci/golangci-lint/internal/x/tools/diff"
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/goformatters"
	"github.com/golangci/golangci-lint/pkg/goformatters/gci"
	"github.com/golangci/golangci-lint/pkg/goformatters/gofmt"
	"github.com/golangci/golangci-lint/pkg/goformatters/gofumpt"
	"github.com/golangci/golangci-lint/pkg/goformatters/goimports"
=======
package processors

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/golangci/golangci-lint/internal/go/robustio"
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/fsutils"
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
	"github.com/golangci/golangci-lint/pkg/timeutils"
)

var _ Processor = (*Fixer)(nil)

<<<<<<< HEAD
const filePerm = 0644

=======
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
type Fixer struct {
	cfg       *config.Config
	log       logutils.Log
	fileCache *fsutils.FileCache
	sw        *timeutils.Stopwatch
<<<<<<< HEAD
	formatter *goformatters.MetaFormatter
}

func NewFixer(cfg *config.Config, log logutils.Log, fileCache *fsutils.FileCache, formatter *goformatters.MetaFormatter) *Fixer {
=======
}

func NewFixer(cfg *config.Config, log logutils.Log, fileCache *fsutils.FileCache) *Fixer {
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	return &Fixer{
		cfg:       cfg,
		log:       log,
		fileCache: fileCache,
		sw:        timeutils.NewStopwatch("fixer", log),
<<<<<<< HEAD
		formatter: formatter,
=======
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	}
}

func (Fixer) Name() string {
	return "fixer"
}

func (p Fixer) Process(issues []result.Issue) ([]result.Issue, error) {
	if !p.cfg.Issues.NeedFix {
		return issues, nil
	}

<<<<<<< HEAD
	p.log.Infof("Applying suggested fixes")

	notFixableIssues, err := timeutils.TrackStage(p.sw, "all", func() ([]result.Issue, error) {
		return p.process(issues)
	})
	if err != nil {
		p.log.Warnf("Failed to fix issues: %v", err)
=======
	outIssues := make([]result.Issue, 0, len(issues))
	issuesToFixPerFile := map[string][]result.Issue{}
	for i := range issues {
		issue := &issues[i]
		if issue.Replacement == nil {
			outIssues = append(outIssues, *issue)
			continue
		}

		issuesToFixPerFile[issue.FilePath()] = append(issuesToFixPerFile[issue.FilePath()], *issue)
	}

	for file, issuesToFix := range issuesToFixPerFile {
		err := p.sw.TrackStageErr("all", func() error {
			return p.fixIssuesInFile(file, issuesToFix)
		})
		if err != nil {
			p.log.Errorf("Failed to fix issues in file %s: %s", file, err)

			// show issues only if can't fix them
			outIssues = append(outIssues, issuesToFix...)
		}
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	}

	p.printStat()

<<<<<<< HEAD
	return notFixableIssues, nil
}

//nolint:funlen,gocyclo // This function should not be split.
func (p Fixer) process(issues []result.Issue) ([]result.Issue, error) {
	// filenames / linters / edits
	editsByLinter := make(map[string]map[string][]diff.Edit)

	formatters := []string{gofumpt.Name, goimports.Name, gofmt.Name, gci.Name}

	var notFixableIssues []result.Issue

	toBeFormattedFiles := make(map[string]struct{})

	for i := range issues {
		issue := issues[i]

		if slices.Contains(formatters, issue.FromLinter) {
			toBeFormattedFiles[issue.FilePath()] = struct{}{}
			continue
		}

		if issue.SuggestedFixes == nil || skipNoTextEdit(&issue) {
			notFixableIssues = append(notFixableIssues, issue)
			continue
		}

		for _, sf := range issue.SuggestedFixes {
			for _, edit := range sf.TextEdits {
				start, end := edit.Pos, edit.End
				if start > end {
					return nil, fmt.Errorf("%q suggests invalid fix: pos (%v) > end (%v)",
						issue.FromLinter, edit.Pos, edit.End)
				}

				edit := diff.Edit{
					Start: int(start),
					End:   int(end),
					New:   string(edit.NewText),
				}

				if _, ok := editsByLinter[issue.FilePath()]; !ok {
					editsByLinter[issue.FilePath()] = make(map[string][]diff.Edit)
				}

				editsByLinter[issue.FilePath()][issue.FromLinter] = append(editsByLinter[issue.FilePath()][issue.FromLinter], edit)
			}
		}
	}

	// Validate and group the edits to each actual file.
	editsByPath := make(map[string][]diff.Edit)
	for path, linterToEdits := range editsByLinter {
		excludedLinters := make(map[string]struct{})

		linters := maps.Keys(linterToEdits)

		// Does any linter create conflicting edits?
		for _, linter := range linters {
			edits := linterToEdits[linter]
			if _, invalid := validateEdits(edits); invalid > 0 {
				name, x, y := linter, edits[invalid-1], edits[invalid]
				excludedLinters[name] = struct{}{}

				err := diff3Conflict(path, name, name, []diff.Edit{x}, []diff.Edit{y})
				// TODO(ldez) TUI?
				p.log.Warnf("Changes related to %q are skipped for the file %q: %v",
					name, path, err)
			}
		}

		// Does any pair of different linters create edits that conflict?
		for j := range linters {
			for k := range linters[:j] {
				x, y := linters[j], linters[k]
				if x > y {
					x, y = y, x
				}

				_, foundX := excludedLinters[x]
				_, foundY := excludedLinters[y]
				if foundX || foundY {
					continue
				}

				xedits, yedits := linterToEdits[x], linterToEdits[y]

				combined := slices.Concat(xedits, yedits)

				if _, invalid := validateEdits(combined); invalid > 0 {
					excludedLinters[x] = struct{}{}
					p.log.Warnf("Changes related to %q are skipped for the file %q due to conflicts with %q.", x, path, y)
				}
			}
		}

		var edits []diff.Edit
		for linter := range linterToEdits {
			if _, found := excludedLinters[linter]; !found {
				edits = append(edits, linterToEdits[linter]...)
			}
		}

		editsByPath[path], _ = validateEdits(edits) // remove duplicates. already validated.
	}

	var editError error

	var formattedFiles []string

	// Now we've got a set of valid edits for each file. Apply them.
	for path, edits := range editsByPath {
		contents, err := p.fileCache.GetFileBytes(path)
		if err != nil {
			editError = errors.Join(editError, fmt.Errorf("%s: %w", path, err))
			continue
		}

		out, err := diff.ApplyBytes(contents, edits)
		if err != nil {
			editError = errors.Join(editError, fmt.Errorf("%s: %w", path, err))
			continue
		}

		// Try to format the file.
		out = p.formatter.Format(path, out)

		if err := os.WriteFile(path, out, filePerm); err != nil {
			editError = errors.Join(editError, fmt.Errorf("%s: %w", path, err))
			continue
		}

		formattedFiles = append(formattedFiles, path)
	}

	for path := range toBeFormattedFiles {
		// Skips files already formatted by the previous fix step.
		if !slices.Contains(formattedFiles, path) {
			content, err := p.fileCache.GetFileBytes(path)
			if err != nil {
				p.log.Warnf("Error reading file %s: %v", path, err)
				continue
			}

			out := p.formatter.Format(path, content)

			if err := os.WriteFile(path, out, filePerm); err != nil {
				editError = errors.Join(editError, fmt.Errorf("%s: %w", path, err))
				continue
			}
		}
	}

	return notFixableIssues, editError
=======
	return outIssues, nil
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
}

func (Fixer) Finish() {}

<<<<<<< HEAD
func (p Fixer) printStat() {
	p.sw.PrintStages()
}

func skipNoTextEdit(issue *result.Issue) bool {
	var onlyMessage int
	for _, sf := range issue.SuggestedFixes {
		if len(sf.TextEdits) == 0 {
			onlyMessage++
		}
	}

	return len(issue.SuggestedFixes) == onlyMessage
}

// validateEdits returns a list of edits that is sorted and
// contains no duplicate edits. Returns the index of some
// overlapping adjacent edits if there is one and <0 if the
// edits are valid.
//
//nolint:gocritic // Copy of go/analysis/internal/checker/checker.go
func validateEdits(edits []diff.Edit) ([]diff.Edit, int) {
	if len(edits) == 0 {
		return nil, -1
	}

	equivalent := func(x, y diff.Edit) bool {
		return x.Start == y.Start && x.End == y.End && x.New == y.New
	}

	diff.SortEdits(edits)

	unique := []diff.Edit{edits[0]}

	invalid := -1

	for i := 1; i < len(edits); i++ {
		prev, cur := edits[i-1], edits[i]
		// We skip over equivalent edits without considering them
		// an error. This handles identical edits coming from the
		// multiple ways of loading a package into a
		// *go/packages.Packages for testing, e.g. packages "p" and "p [p.test]".
		if !equivalent(prev, cur) {
			unique = append(unique, cur)
			if prev.End > cur.Start {
				invalid = i
			}
		}
	}
	return unique, invalid
}

// diff3Conflict returns an error describing two conflicting sets of
// edits on a file at path.
// Copy of go/analysis/internal/checker/checker.go
func diff3Conflict(path, xlabel, ylabel string, xedits, yedits []diff.Edit) error {
	contents, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	oldlabel, old := "base", string(contents)

	xdiff, err := diff.ToUnified(oldlabel, xlabel, old, xedits, diff.DefaultContextLines)
	if err != nil {
		return err
	}
	ydiff, err := diff.ToUnified(oldlabel, ylabel, old, yedits, diff.DefaultContextLines)
	if err != nil {
		return err
	}

	return fmt.Errorf("conflicting edits from %s and %s on %s\nfirst edits:\n%s\nsecond edits:\n%s",
		xlabel, ylabel, path, xdiff, ydiff)
}
=======
func (p Fixer) fixIssuesInFile(filePath string, issues []result.Issue) error {
	// TODO: don't read the whole file into memory: read line by line;
	// can't just use bufio.scanner: it has a line length limit
	origFileData, err := p.fileCache.GetFileBytes(filePath)
	if err != nil {
		return fmt.Errorf("failed to get file bytes for %s: %w", filePath, err)
	}

	origFileLines := bytes.Split(origFileData, []byte("\n"))

	tmpFileName := filepath.Join(filepath.Dir(filePath), fmt.Sprintf(".%s.golangci_fix", filepath.Base(filePath)))

	tmpOutFile, err := os.Create(tmpFileName)
	if err != nil {
		return fmt.Errorf("failed to make file %s: %w", tmpFileName, err)
	}

	// merge multiple issues per line into one issue
	issuesPerLine := map[int][]result.Issue{}
	for i := range issues {
		issue := &issues[i]
		issuesPerLine[issue.Line()] = append(issuesPerLine[issue.Line()], *issue)
	}

	issues = issues[:0] // reuse the same memory
	for line, lineIssues := range issuesPerLine {
		if mergedIssue := p.mergeLineIssues(line, lineIssues, origFileLines); mergedIssue != nil {
			issues = append(issues, *mergedIssue)
		}
	}

	issues = p.findNotIntersectingIssues(issues)

	if err = p.writeFixedFile(origFileLines, issues, tmpOutFile); err != nil {
		tmpOutFile.Close()
		_ = robustio.RemoveAll(tmpOutFile.Name())
		return err
	}

	tmpOutFile.Close()

	if err = robustio.Rename(tmpOutFile.Name(), filePath); err != nil {
		_ = robustio.RemoveAll(tmpOutFile.Name())
		return fmt.Errorf("failed to rename %s -> %s: %w", tmpOutFile.Name(), filePath, err)
	}

	return nil
}

func (p Fixer) mergeLineIssues(lineNum int, lineIssues []result.Issue, origFileLines [][]byte) *result.Issue {
	origLine := origFileLines[lineNum-1] // lineNum is 1-based

	if len(lineIssues) == 1 && lineIssues[0].Replacement.Inline == nil {
		return &lineIssues[0]
	}

	// check issues first
	for ind := range lineIssues {
		li := &lineIssues[ind]

		if li.LineRange != nil {
			p.log.Infof("Line %d has multiple issues but at least one of them is ranged: %#v", lineNum, lineIssues)
			return &lineIssues[0]
		}

		inline := li.Replacement.Inline

		if inline == nil || len(li.Replacement.NewLines) != 0 || li.Replacement.NeedOnlyDelete {
			p.log.Infof("Line %d has multiple issues but at least one of them isn't inline: %#v", lineNum, lineIssues)
			return li
		}

		if inline.StartCol < 0 || inline.Length <= 0 || inline.StartCol+inline.Length > len(origLine) {
			p.log.Warnf("Line %d (%q) has invalid inline fix: %#v, %#v", lineNum, origLine, li, inline)
			return nil
		}
	}

	return p.applyInlineFixes(lineIssues, origLine, lineNum)
}

func (p Fixer) applyInlineFixes(lineIssues []result.Issue, origLine []byte, lineNum int) *result.Issue {
	sort.Slice(lineIssues, func(i, j int) bool {
		return lineIssues[i].Replacement.Inline.StartCol < lineIssues[j].Replacement.Inline.StartCol
	})

	var newLineBuf bytes.Buffer
	newLineBuf.Grow(len(origLine))

	//nolint:misspell // misspelling is intentional
	// example: origLine="it's becouse of them", StartCol=5, Length=7, NewString="because"

	curOrigLinePos := 0
	for i := range lineIssues {
		fix := lineIssues[i].Replacement.Inline
		if fix.StartCol < curOrigLinePos {
			p.log.Warnf("Line %d has multiple intersecting issues: %#v", lineNum, lineIssues)
			return nil
		}

		if curOrigLinePos != fix.StartCol {
			newLineBuf.Write(origLine[curOrigLinePos:fix.StartCol])
		}
		newLineBuf.WriteString(fix.NewString)
		curOrigLinePos = fix.StartCol + fix.Length
	}
	if curOrigLinePos != len(origLine) {
		newLineBuf.Write(origLine[curOrigLinePos:])
	}

	mergedIssue := lineIssues[0] // use text from the first issue (it's not really used)
	mergedIssue.Replacement = &result.Replacement{
		NewLines: []string{newLineBuf.String()},
	}
	return &mergedIssue
}

func (p Fixer) findNotIntersectingIssues(issues []result.Issue) []result.Issue {
	sort.SliceStable(issues, func(i, j int) bool {
		a, b := issues[i], issues[j]
		return a.Line() < b.Line()
	})

	var ret []result.Issue
	var currentEnd int
	for i := range issues {
		issue := &issues[i]
		rng := issue.GetLineRange()
		if rng.From <= currentEnd {
			p.log.Infof("Skip issue %#v: intersects with end %d", issue, currentEnd)
			continue // skip intersecting issue
		}
		p.log.Infof("Fix issue %#v with range %v", issue, issue.GetLineRange())
		ret = append(ret, *issue)
		currentEnd = rng.To
	}

	return ret
}

func (p Fixer) writeFixedFile(origFileLines [][]byte, issues []result.Issue, tmpOutFile *os.File) error {
	// issues aren't intersecting

	nextIssueIndex := 0
	for i := 0; i < len(origFileLines); i++ {
		var outLine string
		var nextIssue *result.Issue
		if nextIssueIndex != len(issues) {
			nextIssue = &issues[nextIssueIndex]
		}

		origFileLineNumber := i + 1
		if nextIssue == nil || origFileLineNumber != nextIssue.GetLineRange().From {
			outLine = string(origFileLines[i])
		} else {
			nextIssueIndex++
			rng := nextIssue.GetLineRange()
			if rng.From > rng.To {
				// Maybe better decision is to skip such issues, re-evaluate if regressed.
				p.log.Warnf("[fixer]: issue line range is probably invalid, fix can be incorrect (from=%d, to=%d, linter=%s)",
					rng.From, rng.To, nextIssue.FromLinter,
				)
			}
			i += rng.To - rng.From
			if nextIssue.Replacement.NeedOnlyDelete {
				continue
			}
			outLine = strings.Join(nextIssue.Replacement.NewLines, "\n")
		}

		if i < len(origFileLines)-1 {
			outLine += "\n"
		}
		if _, err := tmpOutFile.WriteString(outLine); err != nil {
			return fmt.Errorf("failed to write output line: %w", err)
		}
	}

	return nil
}

func (p Fixer) printStat() {
	p.sw.PrintStages()
}
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
