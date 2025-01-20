package linter

import (
	"fmt"

	"golang.org/x/tools/go/packages"

	"github.com/golangci/golangci-lint/pkg/config"
)

const (
	PresetBugs        = "bugs"        // Related to bugs detection.
	PresetComment     = "comment"     // Related to comments analysis.
	PresetComplexity  = "complexity"  // Related to code complexity analysis.
	PresetError       = "error"       // Related to error handling analysis.
	PresetFormatting  = "format"      // Related to code formatting.
	PresetImport      = "import"      // Related to imports analysis.
	PresetMetaLinter  = "metalinter"  // Related to linter that contains multiple rules or multiple linters.
	PresetModule      = "module"      // Related to Go modules analysis.
	PresetPerformance = "performance" // Related to performance.
	PresetSQL         = "sql"         // Related to SQL.
	PresetStyle       = "style"       // Related to coding style.
	PresetTest        = "test"        // Related to the analysis of the code of the tests.
	PresetUnused      = "unused"      // Related to the detection of unused code.
)

// LastLinter nolintlint must be last because it looks at the results of all the previous linters for unused nolint directives.
const LastLinter = "nolintlint"

type DeprecationLevel int

const (
	DeprecationNone DeprecationLevel = iota
	DeprecationWarning
	DeprecationError
)

type Deprecation struct {
	Since       string
	Message     string
	Replacement string
	Level       DeprecationLevel
}

type Config struct {
	Linter           Linter
	EnabledByDefault bool

	LoadMode packages.LoadMode

	InPresets        []string
	AlternativeNames []string

	OriginalURL     string // URL of original (not forked) repo, needed for autogenerated README
	Internal        bool   // Internal linters cannot be disabled (ex: typecheck).
	CanAutoFix      bool
	IsSlow          bool
	DoesChangeTypes bool

	Since       string
	Deprecation *Deprecation
}

func (lc *Config) WithEnabledByDefault() *Config {
	lc.EnabledByDefault = true
	return lc
}

func (lc *Config) WithInternal() *Config {
	lc.Internal = true
	return lc
}

func (lc *Config) ConsiderSlow() *Config {
	lc.IsSlow = true
	return lc
}

func (lc *Config) IsSlowLinter() bool {
	return lc.IsSlow
}

func (lc *Config) WithLoadFiles() *Config {
	lc.LoadMode |= packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedModule
	return lc
}

func (lc *Config) WithLoadForGoAnalysis() *Config {
	lc = lc.WithLoadFiles()
	lc.LoadMode |= packages.NeedImports | packages.NeedDeps | packages.NeedExportFile | packages.NeedTypesSizes
	lc.IsSlow = true
	return lc
}

func (lc *Config) WithPresets(presets ...string) *Config {
	lc.InPresets = presets
	return lc
}

func (lc *Config) WithURL(url string) *Config {
	lc.OriginalURL = url
	return lc
}

func (lc *Config) WithAlternativeNames(names ...string) *Config {
	lc.AlternativeNames = names
	return lc
}

func (lc *Config) WithAutoFix() *Config {
	lc.CanAutoFix = true
	return lc
}

func (lc *Config) WithChangeTypes() *Config {
	lc.DoesChangeTypes = true
	return lc
}

func (lc *Config) WithSince(version string) *Config {
	lc.Since = version
	return lc
}

func (lc *Config) Deprecated(message, version, replacement string, level DeprecationLevel) *Config {
	lc.Deprecation = &Deprecation{
		Since:       version,
		Message:     message,
		Replacement: replacement,
		Level:       level,
	}
	return lc
}

func (lc *Config) DeprecatedWarning(message, version, replacement string) *Config {
	return lc.Deprecated(message, version, replacement, DeprecationWarning)
}

func (lc *Config) DeprecatedError(message, version, replacement string) *Config {
	return lc.Deprecated(message, version, replacement, DeprecationError)
}

func (lc *Config) IsDeprecated() bool {
	return lc.Deprecation != nil
}

func (lc *Config) AllNames() []string {
	return append([]string{lc.Name()}, lc.AlternativeNames...)
}

func (lc *Config) Name() string {
	return lc.Linter.Name()
}

func (lc *Config) WithNoopFallback(cfg *config.Config, cond func(cfg *config.Config) error) *Config {
	if err := cond(cfg); err != nil {
		lc.Linter = NewNoop(lc.Linter, err.Error())
		lc.LoadMode = 0

		return lc.WithLoadFiles()
	}

	return lc
}

func IsGoLowerThanGo122() func(cfg *config.Config) error {
<<<<<<< HEAD
	return isGoLowerThanGo("1.22")
}

func IsGoLowerThanGo124() func(cfg *config.Config) error {
	return isGoLowerThanGo("1.24")
}

func isGoLowerThanGo(v string) func(cfg *config.Config) error {
	return func(cfg *config.Config) error {
		if cfg == nil || config.IsGoGreaterThanOrEqual(cfg.Run.Go, v) {
			return nil
		}

		return fmt.Errorf("this linter is disabled because the Go version (%s) of your project is lower than Go %s", cfg.Run.Go, v)
=======
	return func(cfg *config.Config) error {
		if cfg == nil || config.IsGoGreaterThanOrEqual(cfg.Run.Go, "1.22") {
			return nil
		}

		return fmt.Errorf("this linter is disabled because the Go version (%s) of your project is lower than Go 1.22", cfg.Run.Go)
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	}
}

func NewConfig(linter Linter) *Config {
	lc := &Config{
		Linter: linter,
	}
	return lc.WithLoadFiles()
}
