package main

import (
	"github.com/hailocab/mig31/config"
	"github.com/hailocab/mig31/migration"
	"github.com/hailocab/mig31/runtime"
)

func main() {
	rtConfig := runtime.Flags()

	configErr := rtConfig.Validate()
	if configErr != nil {
		runtime.Usage()
		runtime.ExitWithError(configErr, configErr.ExitCode())
	}

	env, envConfigErr := config.ReadEnvConfig(rtConfig)
	if envConfigErr != nil {
		runtime.ExitWithError(envConfigErr, runtime.ExitErrorReadingEnvConfig)
	}

	if rtConfig.Offline {
		Offline(rtConfig, env)
	}

	if rtConfig.DryRun {
		Dryrun()
	}

	if rtConfig.Verbose {
		VerboseApply()
	}

	if !rtConfig.DryRun && !rtConfig.Offline {
		Apply()
	}
}

// Offline generates the full schema and exits.
func Offline(rtConfig *runtime.Config, env *config.Environment) {
	migReader := migration.NewReader(rtConfig.MigrationsPath)
	migs, readErr := migReader.ReadAllMigrations()
	if readErr != nil {
		runtime.ExitWithError(readErr, runtime.ExitErrorReadingEnvConfig)
	}

	appliedSet := FindAppliedSet(rtConfig, env)

	migs.ApplyEnvironmentStrategy(env)

	sourceSet := migs.SourceSet()
	diff := sourceSet.Diff(appliedSet)

	schema, schemaErr := migs.GenerateSchemaFrom(diff)
	if schemaErr != nil {
		runtime.ExitWithError(schemaErr, runtime.ExitNoEnvironmentDefined)
	}
	runtime.ExitWithMessage(schema, 0)
}

// Dryrun connect to get applied set and output schema.
func Dryrun() {
}

func VerboseApply() {
}

func Apply() {
}
