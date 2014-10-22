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
	migrationSchema  = `
  -- @up

  CREATE KEYSPACE "migrations"
    WITH replication = { 'class': '$${placement_strategy}', $${strategy_options} }
    AND durable_writes = true;

  CREATE TABLE migrations.migrations (
    keyspace_name TEXT PRIMARY KEY,
    ticketNumber INT,
    nextTicketNumber INT,
    migration_ids SET<TEXT>
  ) WITH COMPACT STORAGE AND
  compaction={'class': 'SizeTieredCompactionStrategy'} AND
  compression={'sstable_compression': 'SnappyCompressor'};`
)

// main is the entry point for the application.
func main() {
	var (
		err          error
		env          *config.Environment
		availableSet set.Set
		migs         migration.Migrations
		//lockId      int
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

	// initialise the migration schema
	if flags.Initialise {
		err = Initialise(env)
		if err != nil {
			runtime.ExitWithError(err, runtime.ExitUnableToCreateSchema)
		}
		runtime.ExitWithMessage("Created migration schema.", 0)
	}

	// generate diffSet of available and applied migrations.
	availableSet, err = migration.AvailableSet(flags.MigrationsPath)
	if err != nil {
		runtime.ExitWithError(err, runtime.ExitErrorReadingMigrations)
	}

	cl := dao.New(env.Hosts())
	appliedSet, err := cl.FindAppliedSet(env.Keyspace)
	if err != nil {
		runtime.ExitWithError(err, runtime.ExitErrorReadingMigrations)
	}

	diffSet := availableSet.Diff(appliedSet)

	// read migration files that haven't been applied.
	migReader := migration.NewReader(flags.MigrationsPath, diffSet)
	migs, err = migReader.ReadAll()
	if err != nil {
		runtime.ExitWithError(err, runtime.ExitErrorReadingMigrations)
	}

	//lockId, err =
	iters := []migration.MigrationIter{UpdateSourceSet(env.Keyspace), UpdateStrategy(env)}

	if flags.Offline || flags.Verbose {
		iters = append(iters, PrintUp)
	}

	if !flags.DryRun {
		iters = append(iters, UpdateSchema(cl))
	}

	// TODO: (NF 2014-10-21) This should probably error.
	migs.Apply(iters...)
}

func UpdateSchema(cl dao.MigrationClient) migration.MigrationIter {
	return func(m *migration.Migration) (migration *migration.Migration) {
		migration = m
		return
	}
}

func UpdateSourceSet(ks string) migration.MigrationIter {
	return func(m *migration.Migration) (migration *migration.Migration) {
		update := "\nUPDATE migrations.migrations SET migration_ids = migration_ids + {'" + m.Source + "'} WHERE keyspace_name = '" + ks + "';"
		m.UpMigration = m.UpMigration + update
		migration = m
		return
	}
}

func UpdateStrategy(env *config.Environment) migration.MigrationIter {
	return func(m *migration.Migration) (migration *migration.Migration) {
		up := strings.Replace(m.UpMigration, strategyVariable, env.Strategy(), unlimitedReplace)
		up = strings.Replace(up, optionsVariable, env.Options(), unlimitedReplace)
		m.UpMigration = up
		migration = m
		return
	}

}

func PrintUp(m *migration.Migration) (migration *migration.Migration) {
	fmt.Println(m.UpMigration)
	return
}

func Initialise(env *config.Environment) (err error) {

	return
}
