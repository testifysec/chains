package rules

import (
	"github.com/nunnatsa/ginkgolinter/internal/expression"
	"github.com/nunnatsa/ginkgolinter/internal/expression/matcher"
	"github.com/nunnatsa/ginkgolinter/internal/reports"
	"github.com/nunnatsa/ginkgolinter/types"
)

type HaveLen0 struct{}

func (r *HaveLen0) isApplied(gexp *expression.GomegaExpression, config types.Config) bool {
<<<<<<< HEAD
	return gexp.MatcherTypeIs(matcher.HaveLenZeroMatcherType) && !config.AllowHaveLen0
=======
	return gexp.MatcherTypeIs(matcher.HaveLenZeroMatcherType) && !bool(config.AllowHaveLen0)
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
}

func (r *HaveLen0) Apply(gexp *expression.GomegaExpression, config types.Config, reportBuilder *reports.Builder) bool {
	if !r.isApplied(gexp, config) {
		return false
	}
	gexp.SetMatcherBeEmpty()
	reportBuilder.AddIssue(true, wrongLengthWarningTemplate)
	return true
}
