package cli_template

var SeederMongo string = `
package seeder

import (
	"context"

	"{{.DomainPath}}"
	"fmt"
	r "{{.RepositoryPath}}"
	uc "{{.UsecasePath}}"
	"go.mongodb.org/mongo-driver/mongo"
)

// {{.UcFirstName}}Seeder ...
type {{.UcFirstName}}Seeder struct {
}

// Name returns seeder name
func (*{{.UcFirstName}}Seeder) Name() string {
	return "{{.SmallPluralName}}"
}

// Seed seed {{.SmallName}} document/table
func (*{{.UcFirstName}}Seeder) Seed(ctx context.Context, d *mongo.Database) error {
	{{.SmallName}}Uecase := uc.New{{.UcFirstName}}Usecase(r.New{{.UcFirstName}}Mongo(d))
	fmt.Println({{.SmallName}}Uecase)
	{{.SmallPluralName}} := []domain.{{.UcFirstName}}{}

	for _, {{.SmallName}} := range {{.SmallPluralName}} {
		fmt.Println({{.SmallName}})
	}
	return nil
}
`
