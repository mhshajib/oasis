package config

import (
	"fmt"
	"oasis/pkg/utils"
	"os"

	"github.com/spf13/viper"
)

func GetEnvDefault(key, defVal string) string {
	val, ex := os.LookupEnv(key)
	if !ex {
		return defVal
	}
	return val
}

// Init load configurations from config.yml file
func Init() error {
	viper.SetConfigName("oasis")
	viper.SetConfigType("yml")

	cwd, err := os.Getwd()
	if err != nil {
		os.Exit(1)
	}

	rootPath, _, err := utils.GetGoModFile(cwd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	viper.AddConfigPath(rootPath)
	fmt.Print(cwd)
	// Read the configuration
	err = viper.ReadInConfig()
	if err != nil {
		return err
	}

	initConfig()
	return nil
}

// initConfig laod all configurations
func initConfig() {
	loadPaths()
}
