package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	options "github.com/breathbath/go_utils/utils/config"
	"github.com/breathbath/go_utils/utils/env"
	"github.com/breathbath/go_utils/utils/fs"
	"github.com/sirupsen/logrus"
)

const (
	defaultPath = ".config/rportcli/config.json"
)

func getConfigLocation() (configPath string) {
	configPathFromEnv := env.ReadEnv("CONFIG_PATH", "")
	if configPathFromEnv != "" {
		configPath = configPathFromEnv
		return
	}

	usr, err := user.Current()
	if err != nil {
		logrus.Warnf("failed to read current user data: %v", err)
		configPath = "config.yaml"
		return
	}

	pathParts := []string{usr.HomeDir}
	pathParts = append(pathParts, strings.Split(defaultPath, "/")...)
	configPath = filepath.Join(pathParts...)
	return
}

// GetConfig reads config data from location
func GetConfig() (params *options.ParameterBag, err error) {
	configLocation := getConfigLocation()
	if !fs.FileExists(configLocation) {
		err = fmt.Errorf("config file %s doesn't exist", configLocation)
		return
	}

	f, err := os.Open(configLocation)
	if err != nil {
		err = fmt.Errorf("failed to open the file %s: %v", configLocation, err)
		return
	}

	jvp, err := options.NewJsonValuesProvider(f)
	if err != nil {
		return nil, err
	}

	return &options.ParameterBag{
		BaseValuesProvider: jvp,
	}, nil
}

// GetDefaultConfig creates a config with default values
func GetDefaultConfig() (params *options.ParameterBag) {
	vp := options.NewMapValuesProvider(map[string]interface{}{
		ServerURL: DefaultServerURL,
		Password:  "",
		Login:     "",
	})

	return &options.ParameterBag{BaseValuesProvider: vp}
}

// FromValues creates a parameter bag from provided values
func FromValues(inputParams map[string]string) (params *options.ParameterBag) {
	inputParamsI := map[string]interface{}{}
	for k, v := range inputParams {
		inputParamsI[k] = v
	}
	vp := options.NewMapValuesProvider(inputParamsI)

	return &options.ParameterBag{BaseValuesProvider: vp}
}

// WriteConfig will write config values to file system
func WriteConfig(params *options.ParameterBag) (err error) {
	configLocation := getConfigLocation()

	configDir := filepath.Dir(configLocation)
	err = os.MkdirAll(configDir, 0600)
	if err != nil {
		return err
	}

	fileToWrite, err := os.OpenFile(configLocation, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	err = params.BaseValuesProvider.Dump(fileToWrite)
	if err != nil {
		return err
	}

	return
}
