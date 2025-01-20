package nilnil

import (
	"github.com/Antonboom/nilnil/pkg/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

func New(settings *config.NilNilSettings) *goanalysis.Linter {
	a := analyzer.New()

<<<<<<< HEAD
	cfg := make(map[string]map[string]any)
	if settings != nil {
		cfg[a.Name] = map[string]any{
			"detect-opposite": settings.DetectOpposite,
		}
		if len(settings.CheckedTypes) != 0 {
			cfg[a.Name]["checked-types"] = settings.CheckedTypes
=======
	cfgMap := make(map[string]map[string]any)
	if settings != nil {
		cfgMap[a.Name] = map[string]any{
			"detect-opposite": settings.DetectOpposite,
		}
		if len(settings.CheckedTypes) != 0 {
			cfgMap[a.Name]["checked-types"] = settings.CheckedTypes
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
	).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
