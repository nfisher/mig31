package main

import (
	"testing"
)

func Test_flag_defaults_should_match(t *testing.T) {
	expected := &RuntimeConfig{ConfigPath: "config.xml", Environment: "", DryRun: false, Verbose: false}
	actual := Flags()

	if *expected != *actual {
		t.Fatal("Expected", expected, "but was", actual)
	}
}
