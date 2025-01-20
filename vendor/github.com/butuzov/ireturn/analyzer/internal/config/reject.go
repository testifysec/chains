package config

import "github.com/butuzov/ireturn/analyzer/internal/types"

<<<<<<< HEAD
// rejectConfig specifies a list of interfaces (keywords, patterns and regular expressions)
=======
// rejectConfig specifies a list of interfaces (keywords, patters and regular expressions)
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
// that are rejected by ireturn as valid to return, any non listed interface are allowed.
type rejectConfig struct {
	*defaultConfig
}

func rejectAll(patterns []string) *rejectConfig {
	return &rejectConfig{&defaultConfig{List: patterns}}
}

func (rc *rejectConfig) IsValid(i types.IFace) bool {
	return !rc.Has(i)
}
