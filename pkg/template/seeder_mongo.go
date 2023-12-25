package cli_template

var SeederMongo string = `
package seeder

import (
	"context"

	"{{.ModuleName}}/domain"
	"{{.ModuleName}}/pkg/log"
	r "{{.ModuleName}}/{{.SmallName}}/repository"
	uc "{{.ModuleName}}/{{.SmallName}}/usecase"
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

	{{.SmallPluralName}} := []domain.{{.UcFirstName}}{
		{
			FieldOne:     "Field One",
		},
	}

	for _, {{.SmallName}} := range {{.SmallPluralName}} {
		if err := {{.SmallName}}Uecase.Store(ctx, &{{.SmallName}}); err != nil {
			log.Fatal(err)
		}
	}
	return nil
}
`
