package migration

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func SetupMigrationsDir() (dirPath string, err error) {
	dirPath, err = ioutil.TempDir("", "migrations")
	if err != nil {
		return
	}

	var file *os.File
	file, err = ioutil.TempFile(dirPath, "001_migration")
	if err != nil {
		return
	}

	err = ioutil.WriteFile(file.Name(), []byte(validUpDown), os.ModePerm)
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
	reader := NewReader("cruftlord")
	err := reader.collectMigrationNames()

	if err == nil {
		t.Fatal("An error should have been returned.")
	}
}

func Test_available_migrations_should_fail_if_path_is_a_file(t *testing.T) {
	reader := NewReader("config.xml")
	err := reader.collectMigrationNames()

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

	reader := NewReader(dirPath)
	err := reader.collectMigrationNames()
	if err != nil {
		t.Fatal("An error should not have been returned.", err)
	}

	if len(reader.migrationFiles) < 1 {
		t.Fatal("There should've been at least one file.")
	}
}

func Test_read_all_migrations(t *testing.T) {
	dirPath, setupErr := SetupMigrationsDir()
	if setupErr != nil {
		t.Fatal("Unable to setup test.", setupErr)
	}
	defer TeardownMigrationsDir(dirPath, t)

	reader := NewReader(dirPath)

	migs, err := reader.ReadAllMigrations()
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

	if firstMig.Source == "" {
		t.Fatal("Source should not be empty.")
	}

	if strings.Contains(firstMig.Source, "/") {
		t.Fatal("Source should not include directory paths.")
	}
}
