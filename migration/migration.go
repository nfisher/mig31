package migration

import (
	"errors"
	"github.com/hailocab/mig31/config"
	"github.com/hailocab/mig31/set"
	"path"
	"strings"
)

const (
	upMarker         = "-- @up\n"
	downMarker       = "\n-- @down\n"
	minMarkerSplit   = 2
	strategyVariable = "$${placement_strategy}"
	optionsVariable  = "$${strategy_options}"
	unlimitedReplace = -1
)

type Migrations []Migration

type Migration struct {
	// Source is the the source filename for this migration.
	Source        string
	UpMigration   string
	DownMigration string
}

// ParseMigration is used to verify and split a migration file into up and down migration elements.
func ParseMigration(migrationText string, source string) (migration *Migration, err error) {
	migration = &Migration{}
	upMigration := ""
	downMigration := ""

	split := strings.Split(migrationText, upMarker)
	if len(split) != minMarkerSplit {
		err = errors.New("Unable to find up marker in " + source)
		return
	}
	upMigration = split[1]

	split = strings.Split(upMigration, downMarker)
	if len(split) == minMarkerSplit {
		upMigration = split[0]
		downMigration = split[1]
	}

	migration.UpMigration = upMigration
	migration.DownMigration = downMigration
	migration.Source = path.Base(source)

	return
}

// ApplyEnvironmentValues applies the environments relevant strategy to the migration.
func ApplyEnvironmentValues(migration *Migration, environment *config.Environment) (appliedMigration *Migration, err error) {
	appliedMigration = &Migration{}

	appliedMigration.UpMigration = strings.Replace(migration.UpMigration, strategyVariable, environment.Strategy(), unlimitedReplace)
	appliedMigration.UpMigration = strings.Replace(appliedMigration.UpMigration, optionsVariable, environment.Options(), unlimitedReplace)

	appliedMigration.DownMigration = strings.Replace(migration.DownMigration, strategyVariable, environment.Strategy(), unlimitedReplace)
	appliedMigration.DownMigration = strings.Replace(appliedMigration.DownMigration, optionsVariable, environment.Options(), unlimitedReplace)

	return
}

// TODO: (NF 2014-10-09) Not really happy with this method feels to finicky and prone to allowing the user to misuse the API.
func (migs Migrations) ApplyEnvironmentStrategy(env *config.Environment) (err error) {
	var currentMig *Migration
	for pos, mig := range migs {
		currentMig, err = ApplyEnvironmentValues(&mig, env)
		if err != nil {
			return
		}
		migs[pos] = *currentMig
	}

	return
}

func (migs Migrations) GenerateSchemaFrom(appliedSet set.Set) (schema string, err error) {
	for _, mig := range migs {
		if !appliedSet[mig.Source] {
			schema += mig.UpMigration + "\n"
		}
	}
	return
}
