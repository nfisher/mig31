package dao

import (
	"errors"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/hailocab/mig31/set"
	"regexp"
)

const (
	// Keyspace names are 32 or fewer alpha-numeric characters and underscores, the first of which is an alpha character.
	validKeyspaceName = `^[a-zA-Z][a-zA-Z0-9_]{0,31}$`
)

// CassandraClient
type CassandraClient struct {
	config *gocql.ClusterConfig
}

// NewCassandra will initialise the cluster configuration and return a CassandraClient.
func NewCassandra(hosts []string) (client MigrationClient) {
	config := gocql.NewCluster(hosts...)
	client = &CassandraClient{config: config}

	return
}

// FindAppliedSet will find the currently applied migration ids to compare to the local set available in the local migrations folder.
func (cl *CassandraClient) FindAppliedSet(keyspace string) (appliedSet set.Set, err error) {
	var (
		session *gocql.Session
		rows    []map[string]interface{}
	)

	q := `SELECT migration_ids FROM migrations WHERE keyspace_name=?`

	// change the default keyspace and then create the session.
	err = cl.keyspace(keyspaceName)
	if err != nil {
		return
	}

	session, err = cl.createSession()
	if err != nil {
		return
	}
	defer session.Close()

	appliedSetIter := session.Query(q, keyspace).Iter()

	rows, err = appliedSetIter.SliceMap()
	if err != nil {
		err = errors.New("Unable to retrieve Applied Set: " + err.Error())
		return
	}

	if len(rows) > 1 {
		err = errors.New(fmt.Sprintf("Applied Set has too many entries expected at most 1 but was %v.", len(rows)))
		return
	}

	appliedSet = set.New()
	if len(rows) == 0 {
		return
	}

	firstRow := rows[0]

	rawIds, exists := firstRow[migrationIds]
	if !exists {
		err = errors.New("Applied Set looks corrupt entry exists for keyspace but no migration ids were found.")
		return
	}

	ids, ok := rawIds.([]string)
	if !ok {
		err = errors.New("Applied Set has wrong type for migration ids.")
		return
	}

	for _, id := range ids {
		appliedSet.Add(id)
	}

	err = appliedSetIter.Close()
	if err != nil {
		return
	}

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

// keyspace validates the name of the keyspace and if valid sets it for the current configuration.
func (cl *CassandraClient) keyspace(name string) (err error) {
	var (
		ksMatcher *regexp.Regexp
	)

	ksMatcher, err = regexp.Compile(validKeyspaceName)
	if err != nil {
		return
	}

	if !ksMatcher.MatchString(name) {
		err = errors.New("Keyspace " + name + " is invalid.")
		return
	}

	config := cl.config
	config.Keyspace = name

	return
}
