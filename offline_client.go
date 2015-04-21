package main

import (
	"errors"
	"fmt"
)

// Offline client should be used where its desirable to output schema changes but not connect to C*.
type OfflineClient struct{}

func NewOffline(hosts []string) (client MigrationClient) {
	return &OfflineClient{}
}

// FindAppliedSet returns an empty set as there is no way to know what migrations have been run.
func (cl *OfflineClient) FindAppliedSet(keyspace string) (appliedSet StringSet, err error) {
	return NewStringsSet(), nil
}

func (cl *OfflineClient) Identity(keyspace string) (err error) {
	return errors.New("Are ye daft wee boy you kant get an identity offline.")
}

// CreateSchema will print out the keyspace and table for the migration metadata.
func (cl *OfflineClient) CreateSchema(strategy, option string) (err error) {
	fmt.Println("-- create migrations keyspace")
	fmt.Println(migKeyspace(strategy, option))
	fmt.Println("-- create migrations table")
	fmt.Println(migrationTable)
	return
}
