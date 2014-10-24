package dao

import (
	"github.com/gocql/gocql"
	"github.com/hailocab/mig31/set"
)

type CassandraClient struct {
	config *gocql.ClusterConfig
}

// NewCassandra
func NewCassandra(hosts []string) (client MigrationClient) {
	config := gocql.NewCluster(hosts...)
	client = &CassandraClient{config: config}

	return
}

// FindAppliedSet
func (cl *CassandraClient) FindAppliedSet(keyspace string) (appliedSet set.Set, err error) {
	appliedSet = set.New()
	return
}

// CreateSchema will create a migration schema using CQL3.0.
func (cl *CassandraClient) CreateSchema(strategy, option string) (err error) {
	var (
		session *gocql.Session
	)

	session, err = cl.createSession()
	if err != nil {
		return
	}
	defer session.Close()

	err = session.Query(migKeyspace(strategy, option)).Exec()
	if err != nil {
		return
	}

	err = session.Query(migTable()).Exec()
	if err != nil {
		return
	}

	return
}

// createSession will create a default CQL session.
func (cl *CassandraClient) createSession() (session *gocql.Session, err error) {
	config := cl.config
	session, err = config.CreateSession()
	return
}
