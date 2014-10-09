package config

import (
	"io"
	"io/ioutil"
	"os"
)

// ReadConfig reads the environments configuration file that contains C* placement strategy and options.
func ReadConfig(filename string) (envs *Environments, err error) {
	var (
		contents string
		fd       io.Reader
		rawData  []byte
	)

	fd, err = os.Open(filename)
	if err != nil {
		return
	}

	rawData, err = ioutil.ReadAll(fd)
	if err != nil {
		return
	}

	contents = string(rawData)

	envs, err = UnmarshalConfig(contents)
	if err != nil {
		return
	}

	return
}
