package main

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"
)

type MigrationReader struct {
	dirPath        string
	migrationFiles []string
}

// NewReader will prepare a reader for the specified CQL migration file set.
func NewReader(dirPath string, targetSet Set) (reader *MigrationReader) {
	migrationFiles := make([]string, 0, len(targetSet))
	for k := range targetSet {
		migrationFiles = append(migrationFiles, k)
	}
	sort.Strings(migrationFiles)
	reader = &MigrationReader{dirPath: dirPath, migrationFiles: migrationFiles}
	return
}

// AvailableSet returns a set of the migration filenames found in dirPath.
func AvailableSet(dirPath string) (availableSet Set, err error) {
	var (
		info  os.FileInfo
		files []os.FileInfo
	)

	info, err = os.Stat(dirPath)
	if err != nil {
		return
	}

	if !info.IsDir() {
		err = errors.New("You need to specify a directory for the migration location muppet.")
		return
	}

	if !info.IsDir() {
		err = errors.New("You need to specify a directory for the migration location muppet.")
		return
	}

	files, err = ioutil.ReadDir(dirPath)
	if err != nil {
		return
	}

	availableSet = NewStringsSet()

	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".cql") {
			availableSet.Add(f.Name())
		}
	}

	return
}

// ReadAll reads all of the migrations from the specified directory.
func (reader *MigrationReader) ReadAll() (migrations Migrations, err error) {
	var (
		fd    *os.File
		bytes []byte
		mig   *Migration
	)

	dirPath := reader.dirPath
	migrationFiles := reader.migrationFiles

	for _, filename := range migrationFiles {
		fullPath := path.Join(dirPath, filename)
		fd, err = os.Open(fullPath)
		if err != nil {
			return
		}

		bytes, err = ioutil.ReadAll(fd)
		if err != nil {
			return
		}

		mig, err = ParseMigration(string(bytes), fullPath)
		if err != nil {
			return
		}

		migrations = append(migrations, *mig)
	}
	return
}
