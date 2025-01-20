package dupword

import (
	"strings"

	"github.com/Abirdcfly/dupword"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

<<<<<<< HEAD
func New(settings *config.DupWordSettings) *goanalysis.Linter {
	a := dupword.NewAnalyzer()

	cfg := map[string]map[string]any{}
	if settings != nil {
		cfg[a.Name] = map[string]any{
			"keyword": strings.Join(settings.Keywords, ","),
			"ignore":  strings.Join(settings.Ignore, ","),
=======
func New(setting *config.DupWordSettings) *goanalysis.Linter {
	a := dupword.NewAnalyzer()

	cfgMap := map[string]map[string]any{}
	if setting != nil {
		cfgMap[a.Name] = map[string]any{
			"keyword": strings.Join(setting.Keywords, ","),
			"ignore":  strings.Join(setting.Ignore, ","),
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
		}
	}

	return goanalysis.NewLinter(
		a.Name,
		"checks for duplicate words in the source code",
		[]*analysis.Analyzer{a},
<<<<<<< HEAD
		cfg,
=======
		cfgMap,
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
