package bodyclose

import (
	"github.com/timakin/bodyclose/passes/bodyclose"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

func New() *goanalysis.Linter {
	a := bodyclose.Analyzer

	return goanalysis.NewLinter(
		a.Name,
<<<<<<< HEAD
		a.Doc,
=======
		"checks whether HTTP response body is closed successfully",
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
