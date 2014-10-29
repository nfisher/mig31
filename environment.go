package main

import (
	"encoding/json"
	"encoding/xml"
	"strings"
)

const (
	maxFileBufferSize = 8192
)

type Environments struct {
	Environments []Environment `xml:"environments>environment" json:"environments"`
}

type Placement struct {
	Strategy string `xml:"strategy,attr" json:"strategy"`
	Options  string `xml:"options,attr" json:"options"`
}

type Environment struct {
	Name              string    `xml:"name,attr" json:"name"`
	Host              string    `xml:"host,attr" json:"cluster"`
	Keyspace          string    `xml:"keyspace,attr" json:"keyspace"`
	ConfirmIsOptional bool      `xml:"confirmisoptional,omitempty" json:"confirmisoptional,omitempty"`
	Placement         Placement `xml:"placement" json:"placement"`
}

func NewEnvironment(name, host, strategy, options, keyspace string, confirm bool) (environment *Environment) {
	environment = &Environment{Name: name, Host: host, Placement: Placement{Strategy: strategy, Options: options}, Keyspace: keyspace, ConfirmIsOptional: confirm}
	return
}

type StrategyTarget interface {
	CassandraStrategy(strategy, options string)
}

func UnmarshalConfig(config string, isXml bool) (env *Environments, err error) {
	if isXml {
		return unmarshalXmlConfig(config)
	}

	return unmarshalJsonConfig(config)
}

func unmarshalJsonConfig(config string) (env *Environments, err error) {
	env = new(Environments)
	err = json.Unmarshal([]byte(config), env)
	return
}

func unmarshalXmlConfig(config string) (env *Environments, err error) {
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
