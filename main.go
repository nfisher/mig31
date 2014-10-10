package main

import (
	"github.com/hailocab/mig31/config"
	"github.com/hailocab/mig31/migration"
	"github.com/hailocab/mig31/runtime"
)

func main() {
	rtConfig := runtime.Flags()

	configErr := runtime.ValidateConfig(rtConfig)
	if configErr != nil {
		runtime.Usage()
		runtime.ExitWithError(configErr, configErr.ExitCode())
	}

	env, envConfigErr := config.ReadEnvConfig(rtConfig)
	if envConfigErr != nil {
		runtime.ExitWithError(envConfigErr, runtime.ExitErrorReadingEnvConfig)
	}

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

	if rtConfig.Offline {
		runtime.ExitWithMessage(schema, 0)
	}
}
