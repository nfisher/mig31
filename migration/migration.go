package migration

import (
	"errors"
	"path"
	"strings"
)

const (
	upMarker       = "-- @up\n"
	downMarker     = "\n-- @down\n"
	minMarkerSplit = 2
)

type Migrations []Migration

type MigrationIter func(*Migration) *Migration

// Apply will take the functors and transform the migrations passing the result to the next functor.
func (migs Migrations) Apply(fns ...MigrationIter) {
	for _, m := range migs {
		mp := &m
		for _, fn := range fns {
			t := fn(mp)
			if t != nil {
				mp = t
			}
		}
	}
}

type Migration struct {
	// Source is the the source filename for this migration.
	Source        string
	UpMigration   string
	DownMigration string
}

// New will create a new migration with the supplied source name, up and down migration.
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
