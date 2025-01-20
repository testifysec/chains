package result

import (
	"crypto/md5" //nolint:gosec // for md5 hash
	"fmt"
	"go/token"

<<<<<<< HEAD
	"golang.org/x/tools/go/analysis"
=======
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	"golang.org/x/tools/go/packages"
)

type Range struct {
	From, To int
}

<<<<<<< HEAD
=======
type Replacement struct {
	NeedOnlyDelete bool     // need to delete all lines of the issue without replacement with new lines
	NewLines       []string // if NeedDelete is false it's the replacement lines
	Inline         *InlineFix
}

type InlineFix struct {
	StartCol  int // zero-based
	Length    int // length of chunk to be replaced
	NewString string
}

>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
type Issue struct {
	FromLinter string
	Text       string

	Severity string

	// Source lines of a code with the issue to show
	SourceLines []string

<<<<<<< HEAD
	// Pkg is needed for proper caching of linting results
	Pkg *packages.Package `json:"-"`

	Pos token.Position

	LineRange *Range `json:",omitempty"`

	// HunkPos is used only when golangci-lint is run over a diff
	HunkPos int `json:",omitempty"`

	// If we know how to fix the issue we can provide replacement lines
	SuggestedFixes []analysis.SuggestedFix `json:",omitempty"`

=======
	// If we know how to fix the issue we can provide replacement lines
	Replacement *Replacement

	// Pkg is needed for proper caching of linting results
	Pkg *packages.Package `json:"-"`

	LineRange *Range `json:",omitempty"`

	Pos token.Position

	// HunkPos is used only when golangci-lint is run over a diff
	HunkPos int `json:",omitempty"`

>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	// If we are expecting a nolint (because this is from nolintlint), record the expected linter
	ExpectNoLint         bool
	ExpectedNoLintLinter string
}

func (i *Issue) FilePath() string {
	return i.Pos.Filename
}

func (i *Issue) Line() int {
	return i.Pos.Line
}

func (i *Issue) Column() int {
	return i.Pos.Column
}

func (i *Issue) GetLineRange() Range {
	if i.LineRange == nil {
		return Range{
			From: i.Line(),
			To:   i.Line(),
		}
	}

	if i.LineRange.From == 0 {
		return Range{
			From: i.Line(),
			To:   i.Line(),
		}
	}

	return *i.LineRange
}

func (i *Issue) Description() string {
	return fmt.Sprintf("%s: %s", i.FromLinter, i.Text)
}

func (i *Issue) Fingerprint() string {
	firstLine := ""
	if len(i.SourceLines) > 0 {
		firstLine = i.SourceLines[0]
	}

	hash := md5.New() //nolint:gosec // we don't need a strong hash here
	_, _ = fmt.Fprintf(hash, "%s%s%s", i.Pos.Filename, i.Text, firstLine)

	return fmt.Sprintf("%X", hash.Sum(nil))
}
