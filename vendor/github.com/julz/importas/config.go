package importas

import (
	"errors"
	"fmt"
	"regexp"
<<<<<<< HEAD
	"sync"
)

type Config struct {
	RequiredAlias        aliasList
	Rules                []*Rule
	DisallowUnaliased    bool
	DisallowExtraAliases bool
	muRules              sync.Mutex
}

func (c *Config) CompileRegexp() error {
	c.muRules.Lock()
	defer c.muRules.Unlock()
	if c.Rules != nil {
		return nil
	}
	rules := make([]*Rule, 0, len(c.RequiredAlias))
	for _, aliases := range c.RequiredAlias {
		path, alias := aliases[0], aliases[1]
=======
)

type Config struct {
	RequiredAlias        map[string]string
	Rules                []*Rule
	DisallowUnaliased    bool
	DisallowExtraAliases bool
}

func (c *Config) CompileRegexp() error {
	rules := make([]*Rule, 0, len(c.RequiredAlias))
	for path, alias := range c.RequiredAlias {
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
		reg, err := regexp.Compile(fmt.Sprintf("^%s$", path))
		if err != nil {
			return err
		}

		rules = append(rules, &Rule{
			Regexp: reg,
			Alias:  alias,
		})
	}
<<<<<<< HEAD
=======

>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	c.Rules = rules
	return nil
}

func (c *Config) findRule(path string) *Rule {
<<<<<<< HEAD
	c.muRules.Lock()
	rules := c.Rules
	c.muRules.Unlock()
	for _, rule := range rules {
=======
	for _, rule := range c.Rules {
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
		if rule.Regexp.MatchString(path) {
			return rule
		}
	}

	return nil
}

func (c *Config) AliasFor(path string) (string, bool) {
	rule := c.findRule(path)
	if rule == nil {
		return "", false
	}

	alias, err := rule.aliasFor(path)
	if err != nil {
		return "", false
	}

	return alias, true
}

type Rule struct {
	Alias  string
	Regexp *regexp.Regexp
}

func (r *Rule) aliasFor(path string) (string, error) {
	str := r.Regexp.FindString(path)
	if len(str) > 0 {
		return r.Regexp.ReplaceAllString(str, r.Alias), nil
	}

	return "", errors.New("mismatch rule")
}
