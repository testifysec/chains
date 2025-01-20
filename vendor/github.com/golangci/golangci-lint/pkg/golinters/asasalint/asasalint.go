package asasalint

import (
	"github.com/alingse/asasalint"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/golinters/internal"
)

<<<<<<< HEAD
func New(settings *config.AsasalintSettings) *goanalysis.Linter {
	cfg := asasalint.LinterSetting{}
	if settings != nil {
		cfg.Exclude = settings.Exclude
		cfg.NoBuiltinExclusions = !settings.UseBuiltinExclusions
		cfg.IgnoreTest = settings.IgnoreTest
=======
func New(setting *config.AsasalintSettings) *goanalysis.Linter {
	cfg := asasalint.LinterSetting{}
	if setting != nil {
		cfg.Exclude = setting.Exclude
		cfg.NoBuiltinExclusions = !setting.UseBuiltinExclusions
		cfg.IgnoreTest = setting.IgnoreTest
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	}

	a, err := asasalint.NewAnalyzer(cfg)
	if err != nil {
		internal.LinterLogger.Fatalf("asasalint: create analyzer: %v", err)
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
