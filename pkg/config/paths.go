package config

import (
	"github.com/spf13/viper"
)

type PathStruct struct {
	ModuleName  string
	DomainPath  string
	ServicePath string
	ConfigPath  string
}

var paths PathStruct

func Paths() PathStruct {
	return paths
}

func loadPaths() {
	paths = PathStruct{
		DomainPath:  viper.GetString("paths.domain_path"),
		ServicePath: viper.GetString("paths.service_path"),
		ConfigPath:  viper.GetString("paths.config_path"),
		ModuleName:  GetModuleName(),
	}
}
