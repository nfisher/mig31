package main

import (
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
		env          *Environment
		availableSet StringSet
		migs         Migrations
	)

	// parse and validate flags
	flags := ParseFlags()

	err = flags.Validate()
	if err != nil {
		Usage()
		ExitWithError(err, ExitIncorrectFlag)
	}

	// read environment configuration
	env, err = ReadEnvConfig(flags)
	if err != nil {
		ExitWithError(err, ExitErrorReadingEnvConfig)
	}

	cl := NewClient(env.Hosts(), flags.Username, flags.Password)

	// initialise the migration schema.
	if flags.Initialise {
		err = cl.CreateSchema(env.Strategy(), env.Options())
		if err != nil {
			ExitWithError(err, ExitUnableToCreateSchema)
		}
		return
	}

	// output an identity of a keyspace as a SHA digest.
	if flags.Identity {
		err = cl.Identity(env.Keyspace)
		if err != nil {
			ExitWithError(err, ExitUnableToCreateSchema)
		}
		return
	}

	// generate diffSet of available and applied migrations.
	availableSet, err = AvailableSet(flags.MigrationsPath)
	if err != nil {
		ExitWithError(err, ExitErrorReadingMigrations)
	}

	appliedSet, err := cl.FindAppliedSet(env.Keyspace)
	if err != nil {
		ExitWithError(err, ExitErrorReadingMigrations)
	}

	diffSet := availableSet.Diff(appliedSet)

	missingSet := appliedSet.Diff(availableSet)
	if len(missingSet) != 0 {
		fmt.Println(missingSet)
		ExitWithMessage("This migration set does not match at all.", ExitMigrationMismatch)
	}

	// read migration files that haven't been applied.
	migReader := NewReader(flags.MigrationsPath, diffSet)
	migs, err = migReader.ReadAll()
	if err != nil {
		ExitWithError(err, ExitErrorReadingMigrations)
	}

	iters := []MigrationIter{UpdateSourceSet(env.Keyspace), UpdateStrategy(env)}

	if flags.Offline || flags.Verbose || flags.DryRun {
		iters = append(iters, PrintUp)
	}

	if !flags.DryRun && !flags.Offline {
		ExitWithMessage("Sorry not implemented yet", ExitUnimplemented)
	}

	// TODO: (NF 2014-10-21) This should probably error.
	migs.Apply(iters...)
}

// UpdateSchema will apply the perscribed schema changes to the appropriate keyspace.
func UpdateSchema(cl MigrationClient) MigrationIter {
	return func(m *Migration) (migration *Migration) {
		migration = m
		return
	}
}

// UpdateSourceSet will append the migration source id to the set in the migrations table associated with this keyspace.
func UpdateSourceSet(ks string) MigrationIter {
	return func(m *Migration) (migration *Migration) {
		update := "\nUPDATE migrations.migrations SET migration_ids = migration_ids + {'" + m.Source + "'} WHERE keyspace_name = '" + ks + "';"
		m.UpMigration = m.UpMigration + update
		migration = m
		return
	}
}

// UpdateStrategy will update the strategy for stuff.
func UpdateStrategy(env *Environment) MigrationIter {
	return func(m *Migration) (migration *Migration) {
		up := strings.Replace(m.UpMigration, strategyVariable, env.Strategy(), unlimitedReplace)
		up = strings.Replace(up, optionsVariable, env.Options(), unlimitedReplace)
		m.UpMigration = up
		migration = m
		return
	}

}

// PrintUp will print the up migration.
func PrintUp(m *Migration) (migration *Migration) {
	fmt.Println(m.UpMigration)
	return
}

func Initialise(cl MigrationClient) (err error) {

	return
}
