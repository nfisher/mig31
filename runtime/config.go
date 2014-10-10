package runtime

import (
	"flag"
	"github.com/hailocab/mig31/errors"
)

type Config struct {
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
func (rtConfig *Config) Validate() (err *errors.FatalError) {
	if rtConfig.EnvironmentName == "" {
		err = errors.New("EnvironmentName is required!", ExitIncorrectFlag)
		return
	}

	if rtConfig.ConfigPath == "" {
		err = errors.New("Configuration file path cannot be empty!", ExitIncorrectFlag)
		return
	}

	return
}

func New() (rtConfig *Config) {
	rtConfig = &Config{}
	return
}

// Flags parses all of the command-line flags and returns them as a Config.
func Flags() (rtConfig *Config) {
	rtConfig = New()

	flag.BoolVar(&rtConfig.Offline, "offline", false, "Outputs the full schema without connecting to Cassandra.")
	flag.BoolVar(&rtConfig.DryRun, "dryrun", false, "Dry run and display the changes that would be applied. Implies verbose.")
	flag.BoolVar(&rtConfig.Verbose, "verbose", false, "Enable verbose output.")
	flag.StringVar(&rtConfig.ConfigPath, "config", "config.xml", "Configuration file that specifies placement and strategy options.")
	flag.StringVar(&rtConfig.EnvironmentName, "environment", "", "Name of environment to target.")
	flag.StringVar(&rtConfig.MigrationsPath, "migrations", "./migrations", "Path to the migrations folder.")
	flag.Parse()

	return
}
