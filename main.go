package main

import (
	"flag"
	"fmt"
)

type RuntimeConfig struct {
	ConfigPath  string
	Environment string
	DryRun      bool
	Verbose     bool
}

func Flags() (rtConfig *RuntimeConfig) {
	rtConfig = &RuntimeConfig{}

	flag.BoolVar(&rtConfig.DryRun, "dryrun", false, "Dry run and display the output. Implies verbose.")
	flag.BoolVar(&rtConfig.Verbose, "verbose", false, "Enable verbose output.")
	flag.StringVar(&rtConfig.ConfigPath, "config", "config.xml", "Configuration file that specifies placement and strategy options.")
	flag.StringVar(&rtConfig.Environment, "environment", "", "Environment to target.")
	flag.Parse()

	return
}

func main() {
	rtConfig := Flags()

	fmt.Println("Hello: ", rtConfig.ConfigPath)
}
