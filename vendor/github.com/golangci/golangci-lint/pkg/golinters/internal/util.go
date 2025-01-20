package internal

import (
	"fmt"
<<<<<<< HEAD
=======
	"path/filepath"
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	"strings"

	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
<<<<<<< HEAD
	"github.com/golangci/golangci-lint/pkg/goanalysis"
=======
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
)

func FormatCode(code string, _ *config.Config) string {
	if strings.Contains(code, "`") {
		return code // TODO: properly escape or remove
	}

	return fmt.Sprintf("`%s`", code)
}

<<<<<<< HEAD
func GetGoFileNames(pass *analysis.Pass) []string {
	var filenames []string

	for _, f := range pass.Files {
		position, b := goanalysis.GetGoFilePosition(pass, f)
		if !b {
			continue
		}

		filenames = append(filenames, position.Filename)
	}

	return filenames
=======
func GetFileNames(pass *analysis.Pass) []string {
	var fileNames []string
	for _, f := range pass.Files {
		fileName := pass.Fset.PositionFor(f.Pos(), true).Filename
		ext := filepath.Ext(fileName)
		if ext != "" && ext != ".go" {
			// position has been adjusted to a non-go file, revert to original file
			fileName = pass.Fset.PositionFor(f.Pos(), false).Filename
		}
		fileNames = append(fileNames, fileName)
	}
	return fileNames
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
}
