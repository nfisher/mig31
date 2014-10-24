package dao

import (
	"fmt"
	"github.com/hailocab/mig31/set"
)

type OfflineClient struct{}

func NewOffline(hosts []string) (client MigrationClient) {
	client = &OfflineClient{}
	return
}

// FindAppliedSet find the set of migrations that are currently applied.
func (cl *OfflineClient) FindAppliedSet(keyspace string) (appliedSet set.Set, err error) {
	appliedSet = set.New()
	return
}

func (cl *OfflineClient) Lock() (ticketNum int, err error) {
	ticketNum = 1
	return
}

func (cl *OfflineClient) CreateSchema(strategy, option string) (err error) {
	fmt.Println("-- create migrations keyspace")
	fmt.Println(migKeyspace(strategy, option))
	fmt.Println("-- create migrations table")
	fmt.Println(migrationTable)
	return
}
