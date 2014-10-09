package config

import (
	"encoding/xml"
)

const (
	maxFileBufferSize = 8192
)

type Environments struct {
	Environments []Environment `xml:"environments>environment"`
}

type Placement struct {
	Strategy string `xml:"strategy,attr"`
	Options  string `xml:"options,attr"`
}

type Environment struct {
	Name              string    `xml:"name,attr"`
	Host              string    `xml:"host,attr"`
	ConfirmIsOptional bool      `xml:"confirmisoptional,omitempty"`
	Placement         Placement `xml:"placement"`
}

func NewEnvironment(name, host, strategy, options string) (environment *Environment) {
	environment = &Environment{Name: name, Host: host}
	environment.Placement.Strategy = strategy
	environment.Placement.Options = options
	return
}

type StrategyTarget interface {
	CassandraStrategy(strategy, options string)
}

func UnmarshalConfig(config string) (env *Environments, err error) {
	env = new(Environments)
	err = xml.Unmarshal([]byte(config), env)
	return
}

// Strategy returns the placement strategy for the environment config.
func (environment *Environment) Strategy() (strategy string) {
	return environment.Placement.Strategy
}

// Options returns the strategy options for the environment config.
func (environment *Environment) Options() (strategy string) {
	return environment.Placement.Options
}

func (environment *Environments) Get(name string) (env *Environment) {
	envs := environment.Environments
	for i := 0; i < len(envs); i++ {
		if envs[i].Name == name {
			env = &envs[i]
		}
	}
	return
}

func (environments *Environments) ApplyStrategy(envName string, target StrategyTarget) {
	envs := environments.Environments
	target.CassandraStrategy("boop", "boop")
	for i := 0; i < len(envs); i++ {
		if envs[i].Name == envName {
			target.CassandraStrategy(envs[i].Strategy(), envs[i].Options())
		}
	}
}
