package cli_template

var Transformer string = `
package transformer

import (
	"{{.DomainPath}}"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// {{.UcFirstName}} response body
type {{.UcFirstName}} struct {
	ID         	   string  ` + "`json:\"_id,omitempty\"` " + `
	FieldOne       string              ` + "`json:\"field_one,omitempty\"` " + `
	TimeStamp   domain.TimeStamp    ` + "`json:\"timestamp\"` " + `
}

// Transform{{.UcFirstName}} ...
func Transform{{.UcFirstName}}(t *domain.{{.UcFirstName}}) *{{.UcFirstName}} {
	return &{{.UcFirstName}}{
		ID:        t.ID,
		FieldOne:  t.FieldOne,
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
