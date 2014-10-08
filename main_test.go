package main

import (
	"encoding/xml"
	"testing"
)

func Test_flag_defaults_should_match(t *testing.T) {
	expected := &RuntimeConfig{ConfigPath: "config.xml", Environment: "", DryRun: false, Verbose: false}
	actual := Flags()

	if *expected != *actual {
		t.Fatal("Expected", expected, "but was", actual)
	}
}

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

func Test_config_should_be_deserialised_correctly(t *testing.T) {
	config := `<migrations>
  <environments>
    <environment name="dev" host="10.2.0.21">
      <confirmisoptional>true</confirmisoptional>
      <placement strategy="SimpleStrategy" options="{replication_factor:1}"/>
    </environment>
    <environment name="lve-prem" host="cassandra-premium.eu-west-1a.i.lve.hailocab.net">
      <placement strategy="NetworkTopologyStrategy" options="{eu-west-1 : 3, us-east-1 : 3, ap-northeast-1 : 3}"/>
    </environment>
  </environments>
</migrations>`

	env, err := UnmarshalConfig(config)
	if err != nil {
		t.Fatal("Failed to unmarshal environment config.")
	}

	envCount := len(env.Environments)
	if envCount != 2 {
		t.Fatal("Expected 2 environments but was", envCount)
	}

	devEnv := Environment{Name: "dev", Host: "10.2.0.21", ConfirmIsOptional: true, Placement: Placement{Strategy: "SimpleStrategy", Options: "{replication_factor:1}"}}
	if devEnv != env.Environments[0] {
		t.Fatal("Expected", devEnv, "\nbut was", env.Environments[0])
	}

	lveEnv := Environment{Name: "lve-prem", Host: "cassandra-premium.eu-west-1a.i.lve.hailocab.net", ConfirmIsOptional: false, Placement: Placement{Strategy: "NetworkTopologyStrategy", Options: "{eu-west-1 : 3, us-east-1 : 3, ap-northeast-1 : 3}"}}
	if lveEnv != env.Environments[1] {
		t.Fatal("Expected", lveEnv, "but was", env.Environments[1])
	}
}
