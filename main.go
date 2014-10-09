package main

import (
	"github.com/hailocab/mig31/config"
	"github.com/hailocab/mig31/migration"
)

const (
	exitCodeIncorrectFlag = iota + 1
	exitErrorReadingEnvironmentConfig
	exitNoEnvironmentDefined
)

func main() {
	rtConfig := Flags()

	configErr := ValidateConfig(rtConfig)
	if configErr != nil {
		Usage()
		ExitWithError(configErr, configErr.(*configError).ExitCode())
	}

	envs, envConfigErr := config.ReadConfig(rtConfig.ConfigPath)
	if envConfigErr != nil {
		ExitWithError(envConfigErr, exitErrorReadingEnvironmentConfig)
	}

	env := envs.Get(rtConfig.EnvironmentName)
	if env == nil {
		ExitWithMessage("Environment "+rtConfig.EnvironmentName+" not defined in "+rtConfig.ConfigPath, exitNoEnvironmentDefined)
	}

	migReader := migration.NewReader(rtConfig.MigrationsPath)
	migs, readErr := migReader.ReadAllMigrations()
	if readErr != nil {
		ExitWithError(readErr, exitErrorReadingEnvironmentConfig)
	}

	set := FindAppliedSet(rtConfig, env)

	migs.ApplyEnvironmentStrategy(env)

	schema, schemaErr := migs.GenerateSchemaFrom(set)
	if schemaErr != nil {
		ExitWithError(schemaErr, exitNoEnvironmentDefined)
	}

	if rtConfig.Offline {
		ExitWithMessage(schema, 0)
	}
}
