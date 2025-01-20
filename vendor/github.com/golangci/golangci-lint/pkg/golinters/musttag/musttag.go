package musttag

import (
	"go-simpler.org/musttag"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

<<<<<<< HEAD
func New(settings *config.MustTagSettings) *goanalysis.Linter {
	var funcs []musttag.Func

	if settings != nil {
		for _, fn := range settings.Functions {
=======
func New(setting *config.MustTagSettings) *goanalysis.Linter {
	var funcs []musttag.Func

	if setting != nil {
		for _, fn := range setting.Functions {
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
			funcs = append(funcs, musttag.Func{
				Name:   fn.Name,
				Tag:    fn.Tag,
				ArgPos: fn.ArgPos,
			})
		}
	}

	a := musttag.New(funcs...)

	return goanalysis.
		NewLinter(a.Name, a.Doc, []*analysis.Analyzer{a}, nil).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
