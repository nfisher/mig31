package dao

import (
	"github.com/hailocab/mig31/set"
)

type MigrationClient interface {
	FindAppliedSet(schema string) (appliedSet set.Set, err error)
}

type OfflineClient struct {
}

type CassandraClient struct {
	hosts []string
}

func New(hosts []string) (client MigrationClient) {
	if len(hosts) == 1 && hosts[0] == "" {
		client = &OfflineClient{}
		return
	}

	client = &CassandraClient{hosts: hosts}
	return
}

// FindAppliedSet find the set of migrations that are currently applied.
func (cl *OfflineClient) FindAppliedSet(schema string) (appliedSet set.Set, err error) {
	appliedSet = set.New()
	return
}

func (cl *OfflineClient) Lock() (ticketNum int, err error) {
	ticketNum = 1
	return
}

func (cl *CassandraClient) FindAppliedSet(schema string) (appliedSet set.Set, err error) {
	appliedSet = set.New()
	return
}

func (cl *CassandraClient) Lock() (ticketNum int, err error) {
	ticketNum = 1
	return
}
