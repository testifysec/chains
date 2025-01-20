package config

import "github.com/butuzov/ireturn/analyzer/internal/types"

<<<<<<< HEAD
// allowConfig specifies a list of interfaces (keywords, patterns and regular expressions)
=======
// allowConfig specifies a list of interfaces (keywords, patters and regular expressions)
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
// that are allowed by ireturn as valid to return, any non listed interface are rejected.
type allowConfig struct {
	*defaultConfig
}

func allowAll(patterns []string) *allowConfig {
	return &allowConfig{&defaultConfig{List: patterns}}
}

func (ac *allowConfig) IsValid(i types.IFace) bool {
	return ac.Has(i)
}
