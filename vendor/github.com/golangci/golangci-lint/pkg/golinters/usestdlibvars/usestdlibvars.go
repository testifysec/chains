package usestdlibvars

import (
	"github.com/sashamelentyev/usestdlibvars/pkg/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

<<<<<<< HEAD
func New(settings *config.UseStdlibVarsSettings) *goanalysis.Linter {
	a := analyzer.New()

	cfg := make(map[string]map[string]any)
	if settings != nil {
		cfg[a.Name] = map[string]any{
			analyzer.ConstantKindFlag:       settings.ConstantKind,
			analyzer.CryptoHashFlag:         settings.CryptoHash,
			analyzer.HTTPMethodFlag:         settings.HTTPMethod,
			analyzer.HTTPStatusCodeFlag:     settings.HTTPStatusCode,
			analyzer.OSDevNullFlag:          settings.OSDevNull != nil && *settings.OSDevNull,
			analyzer.RPCDefaultPathFlag:     settings.DefaultRPCPath,
			analyzer.SQLIsolationLevelFlag:  settings.SQLIsolationLevel,
			analyzer.SyslogPriorityFlag:     settings.SyslogPriority != nil && *settings.SyslogPriority,
			analyzer.TimeLayoutFlag:         settings.TimeLayout,
			analyzer.TimeMonthFlag:          settings.TimeMonth,
			analyzer.TimeWeekdayFlag:        settings.TimeWeekday,
			analyzer.TLSSignatureSchemeFlag: settings.TLSSignatureScheme,
=======
func New(cfg *config.UseStdlibVarsSettings) *goanalysis.Linter {
	a := analyzer.New()

	cfgMap := make(map[string]map[string]any)
	if cfg != nil {
		cfgMap[a.Name] = map[string]any{
			analyzer.ConstantKindFlag:       cfg.ConstantKind,
			analyzer.CryptoHashFlag:         cfg.CryptoHash,
			analyzer.HTTPMethodFlag:         cfg.HTTPMethod,
			analyzer.HTTPStatusCodeFlag:     cfg.HTTPStatusCode,
			analyzer.OSDevNullFlag:          cfg.OSDevNull,
			analyzer.RPCDefaultPathFlag:     cfg.DefaultRPCPath,
			analyzer.SQLIsolationLevelFlag:  cfg.SQLIsolationLevel,
			analyzer.SyslogPriorityFlag:     cfg.SyslogPriority,
			analyzer.TimeLayoutFlag:         cfg.TimeLayout,
			analyzer.TimeMonthFlag:          cfg.TimeMonth,
			analyzer.TimeWeekdayFlag:        cfg.TimeWeekday,
			analyzer.TLSSignatureSchemeFlag: cfg.TLSSignatureScheme,
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
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
