package main

import (
	"github.com/hailocab/mig31/config"
	"github.com/hailocab/mig31/migration"
	"github.com/hailocab/mig31/runtime"
)

func main() {
	var (
		err  error
		env  *config.Environment
		migs migration.Migrations
	)

	flags := runtime.ParseFlags()

	err = flags.Validate()
	if err != nil {
		runtime.Usage()
		runtime.ExitWithError(err, runtime.ExitIncorrectFlag)
	}

	env, err = config.ReadEnvConfig(flags)
	if err != nil {
		runtime.ExitWithError(err, runtime.ExitErrorReadingEnvConfig)
	}

	migs, err = ReadMigrations(flags)
	if err != nil {
		runtime.ExitWithError(err, runtime.ExitErrorReadingMigrations)
	}

	if flags.Offline {
		Offline(env, migs)
	}

	if flags.DryRun {
		Dryrun()
	}

	if flags.Verbose {
		VerboseApply()
	}

	if !flags.DryRun && !flags.Offline {
		Apply()
	}
}

func ReadMigrations(flags *runtime.Flags) (migs migration.Migrations, err error) {
	migReader := migration.NewReader(flags.MigrationsPath)
	migs, err = migReader.ReadAll()
	if err != nil {
		return
	}
	return
}

// Offline generates the full schema and exits.
func Offline(env *config.Environment, migs migration.Migrations) {
	appliedSet := FindAppliedSet(env)

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
