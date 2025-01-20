package testpackage

import (
	"strings"

	"github.com/maratori/testpackage/pkg/testpackage"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

<<<<<<< HEAD
func New(settings *config.TestpackageSettings) *goanalysis.Linter {
	a := testpackage.NewAnalyzer()

	var cfg map[string]map[string]any
	if settings != nil {
		cfg = map[string]map[string]any{
			a.Name: {
				testpackage.SkipRegexpFlagName:    settings.SkipRegexp,
				testpackage.AllowPackagesFlagName: strings.Join(settings.AllowPackages, ","),
=======
func New(cfg *config.TestpackageSettings) *goanalysis.Linter {
	a := testpackage.NewAnalyzer()

	var settings map[string]map[string]any
	if cfg != nil {
		settings = map[string]map[string]any{
			a.Name: {
				testpackage.SkipRegexpFlagName:    cfg.SkipRegexp,
				testpackage.AllowPackagesFlagName: strings.Join(cfg.AllowPackages, ","),
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
			},
		}
	}

<<<<<<< HEAD
	return goanalysis.NewLinter(a.Name, a.Doc, []*analysis.Analyzer{a}, cfg).
=======
	return goanalysis.NewLinter(a.Name, a.Doc, []*analysis.Analyzer{a}, settings).
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
		WithLoadMode(goanalysis.LoadModeSyntax)
}
