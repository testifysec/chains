package importas

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

<<<<<<< HEAD
var errWrongAlias = errors.New("import flag must be of form path:alias")

func flags(config *Config) flag.FlagSet {
	fs := flag.FlagSet{}
	fs.Var(&config.RequiredAlias, "alias", "required import alias in form path:alias")
=======
func flags(config *Config) flag.FlagSet {
	fs := flag.FlagSet{}
	fs.Var(stringMap(config.RequiredAlias), "alias", "required import alias in form path:alias")
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	fs.BoolVar(&config.DisallowUnaliased, "no-unaliased", false, "do not allow unaliased imports of aliased packages")
	fs.BoolVar(&config.DisallowExtraAliases, "no-extra-aliases", false, "do not allow non-required aliases")
	return fs
}

<<<<<<< HEAD
type aliasList [][]string

func (v *aliasList) Set(val string) error {
	lastColon := strings.LastIndex(val, ":")
	if lastColon <= 1 {
		return errWrongAlias
	}
	*v = append(*v, []string{val[:lastColon], val[lastColon+1:]})
	return nil
}

func (v *aliasList) String() string {
	return fmt.Sprintf("%v", ([][]string)(*v))
=======
type stringMap map[string]string

func (v stringMap) Set(val string) error {
	spl := strings.SplitN(val, ":", 2)
	if len(spl) != 2 {
		return errors.New("import flag must be of form path:alias")
	}

	v[spl[0]] = spl[1]
	return nil
}

func (v stringMap) String() string {
	return fmt.Sprintf("%v", (map[string]string)(v))
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
}
