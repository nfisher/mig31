package migration

import (
	"errors"
	"github.com/hailocab/mig31/config"
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

func (migs Migrations) ApplyEnvironmentStrategy(env *config.Environment) (err error) {
	var currentMig *Migration
	for i := 0; i < len(migs); i++ {
		currentMig, err = ApplyEnvironmentValues(&migs[i], env)
		if err != nil {
			return
		}
		migs[i] = *currentMig
	}

	return
}
func (migs Migrations) GenerateSchemaFrom(index int) (schema string, err error) {
	for ; index < len(migs); index++ {
		mig := migs[index]
		schema += mig.UpMigration + "\n"
	}
	return
}