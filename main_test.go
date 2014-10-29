package main

import (
	"testing"
)

func Test_flag_defaults_should_match(t *testing.T) {
	expected := NewFlags()
	expected.ConfigPath = "config.json"
	expected.MigrationsPath = "./migrations"
	actual := ParseFlags()

	if *expected != *actual {
		t.Fatal("Expected", expected, "but was", actual)
	}
}

func Test_flag_defaults_should_cause_exit(t *testing.T) {
	//rtConfig := Flags()
	t.Skip("Do something for this test")
}
