package main

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/gocql/gocql"
)

var (
	// Keyspace names are 32 or fewer alpha-numeric characters and underscores, the first of which is an alpha character.
	ksMatcher = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]{0,31}$`)
)

// CassandraClient
type CassandraClient struct {
	config *gocql.ClusterConfig
}

// NewCassandra will initialise the cluster configuration and return a CassandraClient.
func NewCassandra(hosts []string, username, password string) (client MigrationClient) {
	config := gocql.NewCluster(hosts...)
	if username != "" {
		config.Authenticator = &gocql.PasswordAuthenticator{Username: username, Password: password}
	}
	client = &CassandraClient{config: config}

	return client
}

// Identity will output the identity of the specified keyspace as a SHA digest.
func (cl *CassandraClient) Identity(keyspace string) (err error) {
	var (
		session *gocql.Session
		rows    []map[string]interface{}
	)

	err = cl.keyspace("system")
	if err != nil {
		return err
	}

	session, err = cl.createSession()
	if err != nil {
		return err
	}
	defer session.Close()

	q := `SELECT *
          FROM schema_columns
          WHERE keyspace_name = ?`
	colFamilyIter := session.Query(q, keyspace).Iter()
	rows, err = colFamilyIter.SliceMap()
	if err != nil {
		err = errors.New("Unable to retrieve column definition: " + err.Error())
		return
	}

	for _, row := range rows {
		fmt.Println(row)
	}

	return nil
}

// FindAppliedSet will find the currently applied migration ids to compare to the local set available in the local migrations folder.
func (cl *CassandraClient) FindAppliedSet(keyspace string) (appliedSet StringSet, err error) {
	var (
		session *gocql.Session
		rows    []map[string]interface{}
	)

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

	q := `SELECT migration_ids FROM migrations WHERE keyspace_name=?`
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

	appliedSet = NewStringsSet()
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
	if !ksMatcher.MatchString(name) {
		err = errors.New("Keyspace " + name + " is invalid.")
		return
	}

	cl.config.Keyspace = name

	return
}
