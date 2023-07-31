package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/config/secret"
	"github.com/machtwatch/catalystdk/go/log"
	"github.com/spf13/viper"
)

// SetConfig will set config from .env file if it's exist
//
// Otherwise it will set from system's ENV variables.
// Filename should be and env file: .env or .env.* file.
// Dirpath should be in this format: /some/dirpath.
func SetConfig(dirpath string, filename string) {
	filePath := filepath.Join(dirpath, filename)
	fileExist := isFileExist(filePath)

	if fileExist {
		viper.AddConfigPath(dirpath)
		viper.SetConfigFile(filePath)

		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("error reading config file: %+v", err)
		}
	} else {
		viper.AutomaticEnv()
	}

	reloadConfig()
	secret.Reload()
}

// isFileExist check if the file exist on the given file path
func isFileExist(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		return false
	}

	return false
}
