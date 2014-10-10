package config

import (
	"errors"
	"github.com/hailocab/mig31/runtime"
	"io"
	"io/ioutil"
	"os"
)

// ReadEnvConfig reads the environments configuration file that contains C* placement strategy and options.
func ReadEnvConfig(rtConfig *runtime.RuntimeConfig) (env *Environment, err error) {
	var (
		contents string
		fd       io.Reader
		rawData  []byte
		envs     *Environments
	)
	filename := rtConfig.ConfigPath

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

	env = envs.Get(rtConfig.EnvironmentName)
	if env == nil {
		err = errors.New("No Config found for environment" + rtConfig.EnvironmentName + "from" + rtConfig.ConfigPath)
		return
	}

	return
}
