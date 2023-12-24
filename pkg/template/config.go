package cli_template

var Config string = `
package config

import (
	"github.com/spf13/viper"
)

type {{.UcFirstName}}Struct struct {
	FieldOne  string
}

var {{.SmallName}} {{.UcFirstName}}Struct

func {{.UcFirstName}}() {{.UcFirstName}}Struct {
	return {{.SmallName}}
}

func load{{.UcFirstName}}() {
	{{.SmallName}} = {{.UcFirstName}}Struct{
		FieldOne:  viper.GetString("{{.SmallName}}.field_one"),
	}
}
`
