package config

import (
	"github.com/spf13/viper"
)

type PathStruct struct {
	ModuleName    string
	DomainPath    string
	ServicePath   string
	ConfigPath    string
	MigrationPath string
	SeederPath    string
}

var paths PathStruct

func Paths() PathStruct {
	return paths
}

func loadPaths() {
	paths = PathStruct{
		DomainPath:    viper.GetString("paths.domain_path"),
		ServicePath:   viper.GetString("paths.service_path"),
		ConfigPath:    viper.GetString("paths.config_path"),
		MigrationPath: viper.GetString("paths.migration_path"),
		SeederPath:    viper.GetString("paths.seeder_path"),
		ModuleName:    GetModuleName(),
	}
}
