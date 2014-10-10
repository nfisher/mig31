package migration

import (
	"github.com/hailocab/mig31/config"
	"github.com/hailocab/mig31/set"
	"testing"
)

const (
	validUpOnly = `-- @up

create keyspace release with replication = {'class': '$${placement_strategy}', $${strategy_options}};`
	validUpDown = `-- @up

create keyspace release with replication = {'class': '$${placement_strategy}', $${strategy_options}};
-- @down

drop keyspace release;
`
	missingUpMarker = `
create keyspace release with replication = {'class': '$${placement_strategy}', $${strategy_options}};`

	expectedParsedUp = `
create keyspace release with replication = {'class': '$${placement_strategy}', $${strategy_options}};`

	expectedParsedDown = `
drop keyspace release;
`
	expectedAppliedUp = `
create keyspace release with replication = {'class': 'SimpleStrategy', replication_factor: 1};`
)

// end of constants

func Test_should_parse_up_migration_correctly(t *testing.T) {
	mig, err := ParseMigration(validUpOnly, "migrations/001_create_schema.cql")
	if err != nil {
		t.Fatal("Should parse migration successfully.")
	}

	expected := Migration{UpMigration: expectedParsedUp, Source: "001_create_schema.cql"}
	if *mig != expected {
		t.Fatal("Expected", expected, "BUT was", *mig)
	}
}

func Test_should_parse_up_and_down_migration_correctly(t *testing.T) {
	mig, err := ParseMigration(validUpDown, "migrations/001_create_schema.cql")
	if err != nil {
		t.Fatal("Should parse migration successfully.")
	}

	expected := Migration{UpMigration: expectedParsedUp, DownMigration: expectedParsedDown, Source: "001_create_schema.cql"}
	if *mig != expected {
		t.Fatal("Expected <", expected, "> but was <", *mig, ">.")
	}
}

func Test_should_fail_if_up_marker_not_specified(t *testing.T) {
	_, err := ParseMigration(missingUpMarker, "")
	if err == nil {
		t.Fatal("Should fail to parse migration that is missing up marker.")
	}
}

func Test_should_apply_environment_values_to_up_migration(t *testing.T) {
	env := config.NewEnvironment("dev", "localhost", "SimpleStrategy", "replication_factor: 1")
	mig := &Migration{UpMigration: expectedParsedUp}
	appliedMig, err := ApplyEnvironmentValues(mig, env)

	if err != nil {
		t.Fatal("An error occurred applying environment variables.")
	}

	expectedMig := Migration{UpMigration: expectedAppliedUp}

	if *appliedMig != expectedMig {
		t.Fatal("Expected", expectedMig, "but was", *appliedMig)
	}
}

func Test_should_return_full_source_set(t *testing.T) {
	migs := Migrations{Migration{UpMigration: expectedParsedUp, Source: "001_create.cql"}}
	actual := migs.SourceSet()
	expected := set.Set{"001_create.cql": true}

	if !expected.Equal(actual) {
		t.Fatal("Source sets should match.")
	}
}

func Test_should_generate_full_schema_with_empty_set(t *testing.T) {
	migs := Migrations{Migration{UpMigration: expectedParsedUp}}
	_, err := migs.GenerateSchemaFrom(set.New())
	if err != nil {
		t.Fatal(err)
	}
}
