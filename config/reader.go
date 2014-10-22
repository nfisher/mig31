package config

import (
	"errors"
	"github.com/hailocab/mig31/runtime"
	"io"
	"io/ioutil"
	"os"
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
func ReadEnvConfig(flags *runtime.Flags) (env *Environment, err error) {
	var (
		contents string
		envs     *Environments
	)

	contents, err = readFileContents(flags.ConfigPath)
	if err != nil {
		return
	}

	envs, err = UnmarshalConfig(contents)
	if err != nil {
		return
	}

	env = envs.Get(flags.EnvironmentName)
	if env == nil {
		err = errors.New("No Config found for environment" + flags.EnvironmentName + "from" + flags.ConfigPath)
		return
	}

	if flags.Offline {
		optional := env.ConfirmIsOptional
		env = NewEnvironment(env.Name, "", env.Strategy(), env.Options(), env.Keyspace, env.ConfirmIsOptional)
		env.ConfirmIsOptional = optional
	}

	return
}
