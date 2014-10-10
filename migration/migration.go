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

func New(source, up, down string) (migration *Migration) {
	migration = &Migration{Source: source, UpMigration: up, DownMigration: down}
	return
}

// ParseMigration is used to verify and split a migration file into up and down migration elements.
func ParseMigration(migrationText string, source string) (migration *Migration, err error) {
	basename := path.Base(source)
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

	migration = New(basename, upMigration, downMigration)

	return
}

// ApplyEnvironmentValues applies the environments relevant strategy to the migration.
func ApplyEnvironmentValues(migration *Migration, environment *config.Environment) (appliedMigration *Migration, err error) {
	up := strings.Replace(migration.UpMigration, strategyVariable, environment.Strategy(), unlimitedReplace)
	up = strings.Replace(up, optionsVariable, environment.Options(), unlimitedReplace)

	down := strings.Replace(migration.DownMigration, strategyVariable, environment.Strategy(), unlimitedReplace)
	down = strings.Replace(down, optionsVariable, environment.Options(), unlimitedReplace)

	appliedMigration = New(migration.Source, up, down)
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

func (migs Migrations) SourceSet() (sourceSet set.Set) {
	sourceSet = set.New()
	for _, mig := range migs {
		sourceSet.Add(mig.Source)
	}
	return
}

func (migs Migrations) GenerateSchemaFrom(appliedSet set.Set) (schema string, err error) {
	for _, mig := range migs {
		if appliedSet.Contains(mig.Source) {
			schema += mig.UpMigration + "\n"
		}
	}
	return
}
