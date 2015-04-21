package main

import (
	"fmt"
)

const (
	migrationKeyspace = `CREATE KEYSPACE "%v"
    WITH replication = { 'class': '%v', %v }
    AND durable_writes = true;`

	migrationTable = `CREATE TABLE %v.migrations (
    keyspace_name TEXT PRIMARY KEY,
    ticketNumber INT,
    nextTicketNumber INT,
    migration_ids SET<TEXT>
  );`

	keyspaceName = "migrations"
	migrationIds = "migration_ids"
)

type MigrationClient interface {
	FindAppliedSet(keyspace string) (appliedSet StringSet, err error)
	CreateSchema(strategy, options string) (err error)
	Identity(keyspace string) (err error)
}

func NewClient(hosts []string, username, password string) (client MigrationClient) {
	if len(hosts) == 1 && hosts[0] == "" {
		return NewOffline(hosts)
	}

	return NewCassandra(hosts, username, password)
}

func migKeyspace(strategy, options string) (cql string) {
	return fmt.Sprintf(migrationKeyspace, keyspaceName, strategy, options)
}

func migTable() (cql string) {
	return fmt.Sprintf(migrationTable, keyspaceName)
}
