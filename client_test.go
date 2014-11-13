package main

import (
	"testing"
)

func Test_new_with_empty_host_should_return_offline_client(t *testing.T) {
	client := NewClient([]string{""}, "", "")

	_, isOffline := client.(*OfflineClient)

	if !isOffline {
		t.Fatal("Client should be offline but was not.")
	}
}

func Test_new_with_host_should_return_cassandra_client(t *testing.T) {
	client := NewClient([]string{"localhost"}, "", "")

	_, isCassClient := client.(*CassandraClient)

	if !isCassClient {
		t.Fatal("Client should be cassandra but was not.")
	}
}

func Test_new_wit_host_and_password_should_return_password_authenticated_client(t *testing.T) {
	client := NewClient([]string{"localhost"}, "admin", "secret")

	c, isCassClient := client.(*CassandraClient)

	if !isCassClient {
		t.Fatal("Client should be cassandra but was not.")
	}

	if c.config.Authenticator == nil {
		t.Fatal("Client config should have password authentication.")
	}
}

func Test_keyspace_should_error_if_invalid_name(t *testing.T) {
	client := NewClient([]string{"localhost"}, "", "")

	c, _ := client.(*CassandraClient)
	err := c.keyspace("honky *")
	if err == nil {
		t.Fatal("Should have an error for invalid keyspace name.")
	}
}

func Test_keyspace_should_return_nil_error_if_name_is_valid(t *testing.T) {
	client := NewClient([]string{"localhost"}, "", "")

	c, _ := client.(*CassandraClient)
	err := c.keyspace("honky_star")
	if err != nil {
		t.Fatal("Should not have an error for valid keyspace name.")
	}
}
