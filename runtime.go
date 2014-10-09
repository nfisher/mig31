package main

import (
	"flag"
)

type RuntimeConfig struct {
	MigrationsPath  string
	ConfigPath      string
	EnvironmentName string
	DryRun          bool
	Verbose         bool
	Offline         bool
}

func Usage() {
	flag.Usage()
}

// ValidateConfig verifies the provided runtime config is in a sane state returning the first error encountered.
func ValidateConfig(rtConfig *RuntimeConfig) (err error) {
	if rtConfig.EnvironmentName == "" {
		err = newError("EnvironmentName is required!", exitCodeIncorrectFlag)
		return
	}

	if rtConfig.ConfigPath == "" {
		err = newError("Configuration file path cannot be empty!", exitCodeIncorrectFlag)
		return
	}

	return
}

// Flags parses all of the command-line flags and returns them as a RuntimeConfig.
func Flags() (rtConfig *RuntimeConfig) {
	rtConfig = &RuntimeConfig{}

	flag.BoolVar(&rtConfig.Offline, "offline", false, "Outputs the full schema without connecting to Cassandra.")
	flag.BoolVar(&rtConfig.DryRun, "dryrun", false, "Dry run and display the changes that would be applied. Implies verbose.")
	flag.BoolVar(&rtConfig.Verbose, "verbose", false, "Enable verbose output.")
	flag.StringVar(&rtConfig.ConfigPath, "config", "config.xml", "Configuration file that specifies placement and strategy options.")
	flag.StringVar(&rtConfig.EnvironmentName, "environment", "", "Name of environment to target.")
	flag.StringVar(&rtConfig.MigrationsPath, "migrations", "./migrations", "Path to the migrations folder.")
	flag.Parse()

	return
}
