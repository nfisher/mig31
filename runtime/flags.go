package runtime

import (
	"errors"
	"flag"
)

type Flags struct {
	MigrationsPath  string
	ConfigPath      string
	EnvironmentName string
	Username        string
	Password        string
	DryRun          bool
	Verbose         bool
	Offline         bool
	Initialise      bool
	Identity        bool
}

func Usage() {
	flag.Usage()
}

// ValidateConfig verifies the provided runtime config is in a sane state returning the first error encountered.
func (flags *Flags) Validate() (err error) {
	if flags.EnvironmentName == "" {
		err = errors.New("EnvironmentName is required!")
		return
	}

	if flags.ConfigPath == "" {
		err = errors.New("Configuration file path cannot be empty!")
		return
	}

	return
}

func New() (flags *Flags) {
	flags = &Flags{}
	return
}

// Flags parses all of the command-line flags and returns them as a Flags.
func ParseFlags() (flags *Flags) {
	flags = New()

	flag.BoolVar(&flags.Offline, "offline", false, "Outputs the full schema without connecting to Cassandra.")
	flag.BoolVar(&flags.DryRun, "dryrun", false, "Dry run and display the changes that would be applied. Implies verbose.")
	flag.BoolVar(&flags.Verbose, "verbose", false, "Enable verbose output.")
	flag.BoolVar(&flags.Initialise, "init", false, "Initialise the migration keyspace.")
	flag.BoolVar(&flags.Identity, "identity", false, "Describe the keyspace and calculate a SHA as an identity.")
	flag.StringVar(&flags.ConfigPath, "config", "config.xml", "Configuration file that specifies placement and strategy options.")
	flag.StringVar(&flags.EnvironmentName, "env", "", "Name of environment to target.")
	flag.StringVar(&flags.MigrationsPath, "migrations", "./migrations", "Path to the migrations folder.")
	flag.StringVar(&flags.Username, "username", "", "Username to authenticate to Cassandra.")
	flag.StringVar(&flags.Password, "password", "", "Password to authenticate to Cassandra.")
	flag.Parse()

	return
}
