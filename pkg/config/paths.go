package config

import (
	"bufio"
	"clean-cli/pkg/utils"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type PathStruct struct {
	ModuleName string
	DomainPath string
}

var paths PathStruct

func Paths() PathStruct {
	return paths
}

func getModuleName() string {
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

func loadPaths() {
	paths = PathStruct{
		DomainPath: viper.GetString("paths.domain_paths"),
		ModuleName: getModuleName(),
	}
}
