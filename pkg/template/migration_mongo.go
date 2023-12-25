package cli_template

var MigrationMongo string = `
package migration

import (
	"{{.ModuleName}}/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// {{.UcFirstName}} represents {{.SmallName}} migration
type {{.UcFirstName}} struct {
	domain.{{.UcFirstName}}
}

// Name return collection name
func (*{{.UcFirstName}}) Name() string {
	return "{{.SmallPluralName}}"
}

// Indices return collection indices
func (*{{.UcFirstName}}) Indices() []mongo.IndexModel {
	indices := []mongo.IndexModel{
		{
			Keys:    bson.D{{"{{"}}Key: "field_one", Value: "text"{{"}}"}},
			Options: options.Index().SetUnique(true),
		},
	}
	return indices
}
`
