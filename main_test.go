package main

import (
	"github.com/hailocab/mig31/runtime"
	"testing"
)

func Test_flag_defaults_should_match(t *testing.T) {
	expected := runtime.New()
	expected.ConfigPath = "config.xml"
	expected.MigrationsPath = "./migrations"
	actual := runtime.ParseFlags()

	if *expected != *actual {
		t.Fatal("Expected", expected, "but was", actual)
	}
}

func Test_flag_defaults_should_cause_exit(t *testing.T) {
	//rtConfig := Flags()
	t.Skip("Do something for this test")
}
