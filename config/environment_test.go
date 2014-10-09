package config

import (
	"testing"
)

func Test_config_should_be_deserialised_correctly(t *testing.T) {
	config := `<migrations>
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

	lveEnv := Environment{Name: "lve-prem", Host: "lve-prem.local", ConfirmIsOptional: false, Placement: Placement{Strategy: "NetworkTopologyStrategy", Options: "{eu-west-1 : 3, us-east-1 : 3, ap-northeast-1 : 3}"}}
	if lveEnv != env.Environments[1] {
		t.Fatal("Expected", lveEnv, "but was", env.Environments[1])
	}
}
