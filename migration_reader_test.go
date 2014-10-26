package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

const (
	migrationsDirPrefix = "migrations"
	migrationFilePrefix = "001_migration"
)

func SetupMigrationsDir() (dirPath string, err error) {
	dirPath, err = ioutil.TempDir("", migrationsDirPrefix)
	if err != nil {
		return
	}

	var file *os.File
	file, err = ioutil.TempFile(dirPath, migrationFilePrefix)
	if err != nil {
		return
	}

	oldName := file.Name()
	err = ioutil.WriteFile(oldName, []byte(validUpDown), os.ModePerm)
	if err != nil {
		return
	}
	newName := oldName + ".cql"

	// need to rename file so it is picked up by the filter for cql extension.
	err = os.Rename(oldName, newName)
	if err != nil {
		return
	}

	return
}

func TeardownMigrationsDir(dirPath string, t *testing.T) {
	if dirPath == "" {
		return
	}

	err := os.RemoveAll(dirPath)
	if err != nil {
		t.Fatal("Error cleaning up folder", dirPath)
	}
}

func Test_available_migrations_should_fail_if_folder_does_not_exist(t *testing.T) {
	_, err := AvailableSet("cruftlord")

	if err == nil {
		t.Fatal("An error should have been returned.")
	}
}

func Test_available_migrations_should_fail_if_path_is_a_file(t *testing.T) {
	t.Skip("Need to get this working with setup and teardown")
	_, err := AvailableSet("config.xml")

	if err == nil {
		t.Fatal("An error should have been returned.")
	}
}

func Test_available_migrations_should_return_list_of_filenames(t *testing.T) {
	dirPath, setupErr := SetupMigrationsDir()
	if setupErr != nil {
		t.Fatal("Unable to setup test.", setupErr)
	}
	defer TeardownMigrationsDir(dirPath, t)

	s, err := AvailableSet(dirPath)
	if err != nil {
		t.Fatal("An error should not have been returned.", err)
	}

	if len(s) != 1 {
		t.Fatal("There should be exactly one migration.")
	}
}

func Test_read_all_migrations(t *testing.T) {
	dirPath, setupErr := SetupMigrationsDir()
	if setupErr != nil {
		t.Fatal("Unable to setup test.", setupErr)
	}
	defer TeardownMigrationsDir(dirPath, t)

	s, _ := AvailableSet(dirPath)
	reader := NewReader(dirPath, s)

	migs, err := reader.ReadAll()
	if err != nil {
		t.Fatal(err)
	}

	if len(migs) != 1 {
		t.Fatal("There should be exactly one migration.")
	}

	firstMig := migs[0]

	if firstMig.UpMigration == "" {
		t.Fatal("UpMigration should not be empty.")
	}

	if !strings.HasPrefix(firstMig.Source, migrationFilePrefix) {
		t.Fatal("Expected Source to be prefixed by", migrationFilePrefix, "but was", firstMig.Source)
	}

	if strings.Contains(firstMig.Source, "/") {
		t.Fatal("Source should not include directory paths.")
	}
}
