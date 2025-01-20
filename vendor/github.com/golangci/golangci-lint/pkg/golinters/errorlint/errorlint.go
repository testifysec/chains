package errorlint

import (
	"github.com/polyfloyd/go-errorlint/errorlint"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

<<<<<<< HEAD
func New(settings *config.ErrorLintSettings) *goanalysis.Linter {
	var opts []errorlint.Option

	if settings != nil {
		ae := toAllowPairs(settings.AllowedErrors)
=======
func New(cfg *config.ErrorLintSettings) *goanalysis.Linter {
	var opts []errorlint.Option

	if cfg != nil {
		ae := toAllowPairs(cfg.AllowedErrors)
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
		if len(ae) > 0 {
			opts = append(opts, errorlint.WithAllowedErrors(ae))
		}

<<<<<<< HEAD
		aew := toAllowPairs(settings.AllowedErrorsWildcard)
=======
		aew := toAllowPairs(cfg.AllowedErrorsWildcard)
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
		if len(aew) > 0 {
			opts = append(opts, errorlint.WithAllowedWildcard(aew))
		}
	}

	a := errorlint.NewAnalyzer(opts...)

<<<<<<< HEAD
	cfg := map[string]map[string]any{}

	if settings != nil {
		cfg[a.Name] = map[string]any{
			"errorf":       settings.Errorf,
			"errorf-multi": settings.ErrorfMulti,
			"asserts":      settings.Asserts,
			"comparison":   settings.Comparison,
=======
	cfgMap := map[string]map[string]any{}

	if cfg != nil {
		cfgMap[a.Name] = map[string]any{
			"errorf":       cfg.Errorf,
			"errorf-multi": cfg.ErrorfMulti,
			"asserts":      cfg.Asserts,
			"comparison":   cfg.Comparison,
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
		}
	}

	return goanalysis.NewLinter(
		a.Name,
		"errorlint is a linter for that can be used to find code "+
			"that will cause problems with the error wrapping scheme introduced in Go 1.13.",
		[]*analysis.Analyzer{a},
<<<<<<< HEAD
		cfg,
=======
		cfgMap,
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func toAllowPairs(data []config.ErrorLintAllowPair) []errorlint.AllowPair {
	var pairs []errorlint.AllowPair
	for _, allowedError := range data {
		pairs = append(pairs, errorlint.AllowPair(allowedError))
	}
	return pairs
}
