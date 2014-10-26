package main

import (
	"encoding/xml"
	"strings"
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
	Keyspace          string    `xml:"keyspace,attr"`
	ConfirmIsOptional bool      `xml:"confirmisoptional,omitempty"`
	Placement         Placement `xml:"placement"`
}

func NewEnvironment(name, host, strategy, options, keyspace string, confirm bool) (environment *Environment) {
	environment = &Environment{Name: name, Host: host, Placement: Placement{Strategy: strategy, Options: options}, Keyspace: keyspace, ConfirmIsOptional: confirm}
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

func (environment *Environment) Hosts() (hosts []string) {
	hosts = strings.Split(environment.Host, ",")
	return
}

func (environment *Environments) Get(name string) (env *Environment) {
	envs := environment.Environments
	for _, e := range envs {
		if e.Name == name {
			env = &e
			break
		}
	}
	return
}
