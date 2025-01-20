package grouper

import (
	grouper "github.com/leonklingele/grouper/pkg/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

func New(settings *config.GrouperSettings) *goanalysis.Linter {
	a := grouper.New()

<<<<<<< HEAD
	cfg := map[string]map[string]any{}
	if settings != nil {
		cfg[a.Name] = map[string]any{
=======
	linterCfg := map[string]map[string]any{}
	if settings != nil {
		linterCfg[a.Name] = map[string]any{
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
			"const-require-single-const":   settings.ConstRequireSingleConst,
			"const-require-grouping":       settings.ConstRequireGrouping,
			"import-require-single-import": settings.ImportRequireSingleImport,
			"import-require-grouping":      settings.ImportRequireGrouping,
			"type-require-single-type":     settings.TypeRequireSingleType,
			"type-require-grouping":        settings.TypeRequireGrouping,
			"var-require-single-var":       settings.VarRequireSingleVar,
			"var-require-grouping":         settings.VarRequireGrouping,
		}
	}

	return goanalysis.NewLinter(
		a.Name,
		"Analyze expression groups.",
		[]*analysis.Analyzer{a},
<<<<<<< HEAD
		cfg,
=======
		linterCfg,
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
