package gosmopolitan

import (
	"strings"

	"github.com/xen0n/gosmopolitan"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

<<<<<<< HEAD
func New(settings *config.GosmopolitanSettings) *goanalysis.Linter {
	a := gosmopolitan.NewAnalyzer()

	cfg := map[string]map[string]any{}
	if settings != nil {
		cfg[a.Name] = map[string]any{
			"allowtimelocal":  settings.AllowTimeLocal,
			"escapehatches":   strings.Join(settings.EscapeHatches, ","),
			"lookattests":     !settings.IgnoreTests,
			"watchforscripts": strings.Join(settings.WatchForScripts, ","),
=======
func New(s *config.GosmopolitanSettings) *goanalysis.Linter {
	a := gosmopolitan.NewAnalyzer()

	cfgMap := map[string]map[string]any{}
	if s != nil {
		cfgMap[a.Name] = map[string]any{
			"allowtimelocal":  s.AllowTimeLocal,
			"escapehatches":   strings.Join(s.EscapeHatches, ","),
			"lookattests":     !s.IgnoreTests,
			"watchforscripts": strings.Join(s.WatchForScripts, ","),
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
