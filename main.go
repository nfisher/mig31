package main

import (
	"github.com/hailocab/mig31/config"
	"github.com/hailocab/mig31/dao"
	"github.com/hailocab/mig31/migration"
	"github.com/hailocab/mig31/runtime"
	"github.com/hailocab/mig31/set"

	"fmt"
	"strings"
)

const (
	strategyVariable = "$${placement_strategy}"
	optionsVariable  = "$${strategy_options}"
	unlimitedReplace = -1
)

// main is the entry point for the application.
func main() {
	var (
		err          error
		env          *config.Environment
		availableSet set.Set
		migs         migration.Migrations
	)

	// parse and validate flags
	flags := runtime.ParseFlags()

	err = flags.Validate()
	if err != nil {
		runtime.Usage()
		runtime.ExitWithError(err, runtime.ExitIncorrectFlag)
	}

	// read environment configuration
	env, err = config.ReadEnvConfig(flags)
	if err != nil {
		runtime.ExitWithError(err, runtime.ExitErrorReadingEnvConfig)
	}

	cl := dao.New(env.Hosts())

	// initialise the migration schema
	if flags.Initialise {
		err = cl.CreateSchema(env.Strategy(), env.Options())
		if err != nil {
			runtime.ExitWithError(err, runtime.ExitUnableToCreateSchema)
		}
		return
	}

	if flags.Identity {
		err = cl.Identity(env.Keyspace)
		if err != nil {
			runtime.ExitWithError(err, runtime.ExitUnableToCreateSchema)
		}
		return
	}

	// generate diffSet of available and applied migrations.
	availableSet, err = migration.AvailableSet(flags.MigrationsPath)
	if err != nil {
		runtime.ExitWithError(err, runtime.ExitErrorReadingMigrations)
	}

	appliedSet, err := cl.FindAppliedSet(env.Keyspace)
	if err != nil {
		runtime.ExitWithError(err, runtime.ExitErrorReadingMigrations)
	}

	diffSet := availableSet.Diff(appliedSet)

	missingSet := appliedSet.Diff(availableSet)
	if len(missingSet) != 0 {
		fmt.Println(missingSet)
		runtime.ExitWithMessage("This migration set does not match at all.", runtime.ExitMigrationMismatch)
	}

	// read migration files that haven't been applied.
	migReader := migration.NewReader(flags.MigrationsPath, diffSet)
	migs, err = migReader.ReadAll()
	if err != nil {
		runtime.ExitWithError(err, runtime.ExitErrorReadingMigrations)
	}

	iters := []migration.MigrationIter{UpdateSourceSet(env.Keyspace), UpdateStrategy(env)}

	if flags.Offline || flags.Verbose || flags.DryRun {
		iters = append(iters, PrintUp)
	}

	if !flags.DryRun && !flags.Offline {
		runtime.ExitWithMessage("Sorry not implemented yet", runtime.ExitUnimplemented)
	}

	// TODO: (NF 2014-10-21) This should probably error.
	migs.Apply(iters...)
}

// UpdateSchema will apply the perscribed schema changes to the appropriate keyspace.
func UpdateSchema(cl dao.MigrationClient) migration.MigrationIter {
	return func(m *migration.Migration) (migration *migration.Migration) {
		migration = m
		return
	}
}

// UpdateSourceSet will append the migration source id to the set in the migrations table associated with this keyspace.
func UpdateSourceSet(ks string) migration.MigrationIter {
	return func(m *migration.Migration) (migration *migration.Migration) {
		update := "\nUPDATE migrations.migrations SET migration_ids = migration_ids + {'" + m.Source + "'} WHERE keyspace_name = '" + ks + "';"
		m.UpMigration = m.UpMigration + update
		migration = m
		return
	}
}

// UpdateStrategy will update the strategy for stuff.
func UpdateStrategy(env *config.Environment) migration.MigrationIter {
	return func(m *migration.Migration) (migration *migration.Migration) {
		up := strings.Replace(m.UpMigration, strategyVariable, env.Strategy(), unlimitedReplace)
		up = strings.Replace(up, optionsVariable, env.Options(), unlimitedReplace)
		m.UpMigration = up
		migration = m
		return
	}

}

// PrintUp will print the up migration.
func PrintUp(m *migration.Migration) (migration *migration.Migration) {
	fmt.Println(m.UpMigration)
	return
}

func Initialise(cl dao.MigrationClient) (err error) {

	return
}
