package util

import (
	"strings"
	"swiftbeaver-infra/config"
)

type Affixer struct {
	prefix  string
	postfix string

	resourceDelimiter string
	outputDelimiter   string
}

func (a *Affixer) join(delimiter string, parts ...string) string {
	existing := []string{}

	for _, part := range parts {
		if part != "" {
			existing = append(existing, part)
		}
	}

	return strings.Join(existing, delimiter)
}

func (a *Affixer) joinResource(parts ...string) string {
	return a.join(a.resourceDelimiter, parts...)
}

func (a *Affixer) joinOutput(parts ...string) string {
	return a.join(a.outputDelimiter, parts...)
}

func (a *Affixer) addResourcePrefix(name string) string {
	return a.joinResource(a.prefix, name)
}

func (a *Affixer) addResourcePostfix(name string) string {
	return a.joinResource(name, a.postfix)
}

func (a *Affixer) addOutputPrefix(name string) string {
	return a.joinOutput(a.prefix, name)
}

func (a *Affixer) PrefixResource(name string) string {
	return a.addResourcePrefix(name)
}

func (a *Affixer) PostfixResource(name string) string {
	return a.addResourcePostfix(name)
}

func (a *Affixer) AffixResource(name string) string {
	return a.addResourcePostfix(a.addResourcePrefix(name))
}

func (a *Affixer) PrefixOutput(name string) string {
	return a.addOutputPrefix(name)
}

func NewAffixer(cfg *config.Config) *Affixer {
	var postfix string

	if cfg.Environment != config.GlobalEnvironment {
		postfix = string(cfg.Environment)
	}

	return &Affixer{
		prefix:  cfg.ProjectName,
		postfix: postfix,

		resourceDelimiter: "-",
		outputDelimiter:   "::",
	}
}
