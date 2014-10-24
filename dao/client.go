package dao

import (
	"fmt"
	"github.com/hailocab/mig31/set"
)

const (
	migrationKeyspace = `CREATE KEYSPACE "migration"
    WITH replication = { 'class': '%v', %v }
    AND durable_writes = true;`

	migrationTable = `CREATE TABLE migration.migrations (
    keyspace_name TEXT PRIMARY KEY,
    ticketNumber INT,
    nextTicketNumber INT,
    migration_ids SET<TEXT>
  );`
)

type MigrationClient interface {
	FindAppliedSet(keyspace string) (appliedSet set.Set, err error)
	CreateSchema(strategy, options string) (err error)
}

func New(hosts []string) (client MigrationClient) {
	if len(hosts) == 1 && hosts[0] == "" {
		client = NewOffline(hosts)
		return
	}

	client = NewCassandra(hosts)
	return
}

func migKeyspace(strategy, options string) (cql string) {
	cql = fmt.Sprintf(migrationKeyspace, strategy, options)
	return
}

func migTable() (cql string) {
	cql = migrationTable
	return
}
