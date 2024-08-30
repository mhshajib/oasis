package cli_template

var Transformer string = `
package transformer

import (
	"{{.DomainPath}}"
)

// {{.UcFirstName}} response body
type {{.UcFirstName}} struct {
	ID         	   string  ` + "`json:\"_id,omitempty\"` " + ` {{range .Fields}}
    {{.Name}}    {{.Type}}          ` + "`json:\"{{.JsonTag}}\"`" + ` {{end}}
	TimeStamp   domain.TimeStamp    ` + "`json:\"timestamp\"` " + `
}

// Transform{{.UcFirstName}} ...
func Transform{{.UcFirstName}}(t *domain.{{.UcFirstName}}) *{{.UcFirstName}} {
	return &{{.UcFirstName}}{
		ID:        t.ID, {{range .Fields}}
    	{{.Name}}: t.{{.Name}}, {{end}}
		TimeStamp: t.TimeStamp,
	}
}

// Transform{{.UcFirstName}}List ...
func Transform{{.UcFirstName}}List(tl []*domain.{{.UcFirstName}}) []*{{.UcFirstName}} {
	{{.SmallPluralName}} := make([]*{{.UcFirstName}}, 0)
	for _, t := range tl {
		{{.SmallPluralName}} = append({{.SmallPluralName}}, Transform{{.UcFirstName}}(t))
	}
	return {{.SmallPluralName}}
}
`
