package maintidx

import (
	"github.com/yagipy/maintidx"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

<<<<<<< HEAD
func New(settings *config.MaintIdxSettings) *goanalysis.Linter {
	analyzer := maintidx.Analyzer

	cfg := map[string]map[string]any{
		analyzer.Name: {"under": 20},
	}

	if settings != nil {
		cfg[analyzer.Name] = map[string]any{
			"under": settings.Under,
=======
func New(cfg *config.MaintIdxSettings) *goanalysis.Linter {
	analyzer := maintidx.Analyzer

	cfgMap := map[string]map[string]any{
		analyzer.Name: {"under": 20},
	}

	if cfg != nil {
		cfgMap[analyzer.Name] = map[string]any{
			"under": cfg.Under,
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
		}
	}

	return goanalysis.NewLinter(
		analyzer.Name,
		analyzer.Doc,
		[]*analysis.Analyzer{analyzer},
<<<<<<< HEAD
		cfg,
=======
		cfgMap,
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
