package config

import (
	"testing"
)

const (
	config = `<migrations>
  <environments>
    <environment name="dev" host="10.2.0.21">
      <confirmisoptional>true</confirmisoptional>
      <placement strategy="SimpleStrategy" options="{replication_factor:1}"/>
    </environment>
    <environment name="lve-prem" host="lve-prem.local">
      <placement strategy="NetworkTopologyStrategy" options="{eu-west-1 : 3, us-east-1 : 3, ap-northeast-1 : 3}"/>
    </environment>
  </environments>
</migrations>`
)

func Test_config_should_be_deserialised_correctly(t *testing.T) {
	envs, err := UnmarshalConfig(config)
	if err != nil {
		t.Fatal("Failed to unmarshal environment config.")
	}

	envCount := len(envs.Environments)
	if envCount != 2 {
		t.Fatal("Expected 2 environments but was", envCount)
	}

	devEnv := Environment{Name: "dev", Host: "10.2.0.21", ConfirmIsOptional: true, Placement: Placement{Strategy: "SimpleStrategy", Options: "{replication_factor:1}"}}
	if devEnv != envs.Environments[0] {
		t.Fatal("Expected", devEnv, "\nbut was", envs.Environments[0])
	}

	lveEnv := Environment{Name: "lve-prem", Host: "lve-prem.local", ConfirmIsOptional: false, Placement: Placement{Strategy: "NetworkTopologyStrategy", Options: "{eu-west-1 : 3, us-east-1 : 3, ap-northeast-1 : 3}"}}
	if lveEnv != envs.Environments[1] {
		t.Fatal("Expected", lveEnv, "but was", envs.Environments[1])
	}
}

type StrategyStub struct {
	Option    string
	Strategy  string
	CallCount int
	Test      *testing.T
}

func (ss StrategyStub) CassandraStrategy(strategy, options string) {
	ss.CallCount++
	ss.Option = options
	ss.Strategy = strategy
}

func Test_get_should_retrieve_the_correct_environment(t *testing.T) {
	envs, _ := UnmarshalConfig(config)
	expected := Environment{Name: "dev", Host: "10.2.0.21", ConfirmIsOptional: true, Placement: Placement{Strategy: "SimpleStrategy", Options: "{replication_factor:1}"}}
	actual := envs.Get("dev")

	if expected != *actual {
		t.Fatal("Expected", expected, "but was", *actual)
	}
}
