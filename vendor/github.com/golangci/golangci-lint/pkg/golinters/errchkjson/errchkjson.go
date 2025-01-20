package errchkjson

import (
	"github.com/breml/errchkjson"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

<<<<<<< HEAD
func New(settings *config.ErrChkJSONSettings) *goanalysis.Linter {
	a := errchkjson.NewAnalyzer()

	cfg := map[string]map[string]any{}
	cfg[a.Name] = map[string]any{
		"omit-safe": true,
	}
	if settings != nil {
		cfg[a.Name] = map[string]any{
			"omit-safe":          !settings.CheckErrorFreeEncoding,
			"report-no-exported": settings.ReportNoExported,
=======
func New(cfg *config.ErrChkJSONSettings) *goanalysis.Linter {
	a := errchkjson.NewAnalyzer()

	cfgMap := map[string]map[string]any{}
	cfgMap[a.Name] = map[string]any{
		"omit-safe": true,
	}
	if cfg != nil {
		cfgMap[a.Name] = map[string]any{
			"omit-safe":          !cfg.CheckErrorFreeEncoding,
			"report-no-exported": cfg.ReportNoExported,
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
		}
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
<<<<<<< HEAD
		cfg,
=======
		cfgMap,
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
