package migration

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"sort"
)

type MigrationReader struct {
	dirPath        string
	migrationFiles []string
}

func NewReader(dirPath string) (reader *MigrationReader) {
	reader = &MigrationReader{dirPath: dirPath}
	return
}

// collectMigrationNames gives a list of available migration files in ascending order.
func (reader *MigrationReader) collectMigrationNames() (err error) {
	var (
		info       os.FileInfo
		files      []os.FileInfo
		migrations []string
	)
	dirPath := reader.dirPath

	info, err = os.Stat(dirPath)
	if err != nil {
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

	// TODO: (NF 2014-10-09) Filter only CQL files.
	for i := 0; i < len(files); i++ {
		f := files[i]
		if !f.IsDir() {
			migrations = append(migrations, f.Name())
		}
	}

	sort.Strings(migrations)

	reader.migrationFiles = migrations
	return
}

// ReadAllMigrations reads all of the migrations from the specified directory.
func (reader *MigrationReader) ReadAllMigrations() (migrations Migrations, err error) {
	var (
		fd    *os.File
		bytes []byte
		mig   *Migration
	)

	err = reader.collectMigrationNames()
	if err != nil {
		return
	}

	dirPath := reader.dirPath
	migrationFiles := reader.migrationFiles

	for i := 0; i < len(migrationFiles); i++ {
		fullPath := path.Join(dirPath, migrationFiles[i])
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
