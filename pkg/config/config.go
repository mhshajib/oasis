package config

import (
	"bufio"
	"fmt"
	"oasis/pkg/utils"
	"os"
	"strings"

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

func GetModuleName() string {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		os.Exit(1)
	}

	_, goModPath, err := utils.GetGoModFile(cwd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	file, err := os.Open(goModPath)
	if err != nil {
		fmt.Println("Error opening go.mod file:", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "module ") {
			moduleName := strings.TrimSpace(line[len("module "):])
			return moduleName
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	return ""
}

// initConfig laod all configurations
func initConfig() {
	loadPaths()
}
