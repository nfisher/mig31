package config

import (
	"encoding/xml"
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

func UnmarshalConfig(config string) (env *Environments, err error) {
	env = &Environments{}
	err = xml.Unmarshal([]byte(config), env)
	return
}
