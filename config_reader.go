package main

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func readFileContents(filename string) (contents string, err error) {
	var fd io.Reader
	var rawData []byte

	fd, err = os.Open(filename)
	if err != nil {
		return
	}

	rawData, err = ioutil.ReadAll(fd)
	if err != nil {
		return
	}

	contents = string(rawData)
	return
}

// ReadEnvConfig reads the environments configuration file that contains C* host, placement strategy and options.
func ReadEnvConfig(flags *Flags) (env *Environment, err error) {
	var (
		contents string
		envs     *Environments
	)

	contents, err = readFileContents(flags.ConfigPath)
	if err != nil {
		return
	}

	isXml := strings.HasSuffix(flags.ConfigPath, ".xml")

	envs, err = UnmarshalConfig(contents, isXml)
	if err != nil {
		return
	}

	env = envs.Get(flags.EnvironmentName)
	if env == nil {
		err = errors.New("No Config found for environment" + flags.EnvironmentName + "from" + flags.ConfigPath)
		return
	}

	if env.Keyspace == "" {
		err = errors.New("No Keyspace found for environment" + flags.EnvironmentName + "from" + flags.ConfigPath)
		return
	}

	if flags.Offline {
		optional := env.ConfirmIsOptional
		env = NewEnvironment(env.Name, "", env.Strategy(), env.Options(), env.Keyspace, env.ConfirmIsOptional)
		env.ConfirmIsOptional = optional
	}

	return
}
