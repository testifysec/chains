package ginkgolinter

import (
	"github.com/nunnatsa/ginkgolinter"
	"github.com/nunnatsa/ginkgolinter/types"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

func New(settings *config.GinkgoLinterSettings) *goanalysis.Linter {
	cfg := &types.Config{}

	if settings != nil {
		cfg = &types.Config{
<<<<<<< HEAD
			SuppressLen:            settings.SuppressLenAssertion,
			SuppressNil:            settings.SuppressNilAssertion,
			SuppressErr:            settings.SuppressErrAssertion,
			SuppressCompare:        settings.SuppressCompareAssertion,
			SuppressAsync:          settings.SuppressAsyncAssertion,
			ForbidFocus:            settings.ForbidFocusContainer,
			SuppressTypeCompare:    settings.SuppressTypeCompareWarning,
			AllowHaveLen0:          settings.AllowHaveLenZero,
			ForceExpectTo:          settings.ForceExpectTo,
			ValidateAsyncIntervals: settings.ValidateAsyncIntervals,
			ForbidSpecPollution:    settings.ForbidSpecPollution,
			ForceSucceedForFuncs:   settings.ForceSucceedForFuncs,
=======
			SuppressLen:            types.Boolean(settings.SuppressLenAssertion),
			SuppressNil:            types.Boolean(settings.SuppressNilAssertion),
			SuppressErr:            types.Boolean(settings.SuppressErrAssertion),
			SuppressCompare:        types.Boolean(settings.SuppressCompareAssertion),
			SuppressAsync:          types.Boolean(settings.SuppressAsyncAssertion),
			ForbidFocus:            types.Boolean(settings.ForbidFocusContainer),
			SuppressTypeCompare:    types.Boolean(settings.SuppressTypeCompareWarning),
			AllowHaveLen0:          types.Boolean(settings.AllowHaveLenZero),
			ForceExpectTo:          types.Boolean(settings.ForceExpectTo),
			ValidateAsyncIntervals: types.Boolean(settings.ValidateAsyncIntervals),
			ForbidSpecPollution:    types.Boolean(settings.ForbidSpecPollution),
			ForceSucceedForFuncs:   types.Boolean(settings.ForceSucceedForFuncs),
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
		}
	}

	a := ginkgolinter.NewAnalyzerWithConfig(cfg)

	return goanalysis.NewLinter(
		a.Name,
		"enforces standards of using ginkgo and gomega",
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
