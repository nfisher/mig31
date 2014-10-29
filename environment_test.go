package main

import (
	"testing"
)

const (
	xmlConfig = `<migrations>
  <environments>
    <environment name="dev" host="10.2.0.21" keyspace="release">
      <confirmisoptional>true</confirmisoptional>
      <placement strategy="SimpleStrategy" options="{replication_factor:1}"/>
    </environment>
    <environment name="lve-prem" host="lve-prem.local" keyspace="release">
      <placement strategy="NetworkTopologyStrategy" options="{eu-west-1 : 3, us-east-1 : 3, ap-northeast-1 : 3}"/>
    </environment>
  </environments>
</migrations>`

	jsonConfig = `{
  "environments": [
      { 
        "name": "dev", 
        "cluster":"10.2.0.21",
        "confirmisoptional": true,
        "keyspace": "release",
        "placement": { "strategy": "SimpleStrategy", "options": "'replication_factor': 1" }
      },
      { 
        "name": "lve-prem", 
        "cluster": "lve-prem.local",
        "keyspace": "release",
        "placement": { "strategy": "NetworkTopologyStrategy", "options": "eu-west-1 : 3, us-east-1 : 3, ap-northeast-1 : 3" }
      }
    ]
}`
)

func Test_json_config_should_be_deserialised_correctly(t *testing.T) {
	envs, err := unmarshalJsonConfig(jsonConfig)
	if err != nil {
		t.Fatal("Failed to unmarshal environment config.")
	}

	envCount := len(envs.Environments)
	if envCount != 2 {
		t.Fatal("Expected 2 environments but was", envCount)
	}

	devEnv := *NewEnvironment("dev", "10.2.0.21", "SimpleStrategy", "'replication_factor': 1", "release", true)
	actual := envs.Environments[0]
	if devEnv != actual {
		t.Fatal("Expected", devEnv, "\nbut was", envs.Environments[0])
	}

	lveEnv := *NewEnvironment("lve-prem", "lve-prem.local", "NetworkTopologyStrategy", "eu-west-1 : 3, us-east-1 : 3, ap-northeast-1 : 3", "release", false)
	actual = envs.Environments[1]
	if lveEnv != actual {
		t.Fatal("Expected", lveEnv, "but was", envs.Environments[1])
	}
}

func Test_xml_config_should_be_deserialised_correctly(t *testing.T) {
	envs, err := unmarshalXmlConfig(xmlConfig)
	if err != nil {
		t.Fatal("Failed to unmarshal environment config.")
	}

	envCount := len(envs.Environments)
	if envCount != 2 {
		t.Fatal("Expected 2 environments but was", envCount)
	}

	devEnv := *NewEnvironment("dev", "10.2.0.21", "SimpleStrategy", "{replication_factor:1}", "release", true)
	actual := envs.Environments[0]
	if devEnv != actual {
		t.Fatal("Expected", devEnv, "\nbut was", envs.Environments[0])
	}

	lveEnv := *NewEnvironment("lve-prem", "lve-prem.local", "NetworkTopologyStrategy", "{eu-west-1 : 3, us-east-1 : 3, ap-northeast-1 : 3}", "release", false)
	actual = envs.Environments[1]
	if lveEnv != actual {
		t.Fatal("Expected", lveEnv, "but was", envs.Environments[1])
	}
}

func Test_get_should_retrieve_the_correct_environment(t *testing.T) {
	envs, _ := unmarshalXmlConfig(xmlConfig)
	expected := NewEnvironment("dev", "10.2.0.21", "SimpleStrategy", "{replication_factor:1}", "release", true)
	actual := envs.Get("dev")

	if *expected != *actual {
		t.Fatal("Expected", expected, "but was", actual)
	}
}

func Test_hosts_should_return_correct_value_for_single_host_entry(t *testing.T) {
	envs, _ := unmarshalXmlConfig(xmlConfig)
	expected := "10.2.0.21"
	devEnv := envs.Get("dev")
	actual := devEnv.Hosts()

	if len(actual) != 1 {
		t.Fatal("Expected exactly 1 element but was ", len(actual))
	}

	if actual[0] != expected {
		t.Fatal("Expected", expected, "but was", actual[0])
	}
}

func Test_hosts_should_return_empty_value_for_offline_host(t *testing.T) {
	env := NewEnvironment("dev", "", "SimpleStrategy", "{replication_factor:1}", "release", true)
	actual := env.Hosts()
	expected := ""

	if len(actual) != 1 {
		t.Fatal("Expected 1 element but was ", len(actual))
	}

	if actual[0] != expected {
		t.Fatal("Expected", expected, "but was", actual[0])
	}
}
