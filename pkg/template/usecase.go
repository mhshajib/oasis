package cli_template

var UseCase string = `
package usecase

import (
	"context"

	"{{.ModuleName}}/domain"
)

// {{.UcFirstName}}Usecase represents {{.SmallName}} usecases
type {{.UcFirstName}}Usecase struct {
	{{.SmallName}}Repository domain.{{.UcFirstName}}Repository
}

// New{{.UcFirstName}}Usecase return {{.SmallName}} usecase instance
func New{{.UcFirstName}}Usecase(r domain.{{.UcFirstName}}Repository) *{{.UcFirstName}}Usecase {
	return &{{.UcFirstName}}Usecase{
		{{.SmallName}}Repository: r,
	}
}

// Store insert a new {{.SmallName}} to storage
func (u *{{.UcFirstName}}Usecase) Store(ctx context.Context, {{.SmallName}} *domain.{{.UcFirstName}}) error {
	{{.SmallName}}.PopulateCreateTimeStamp() // generate created_at & updated_at field

	return u.{{.SmallName}}Repository.Store(ctx, {{.SmallName}})
}

// Fetch list {{.SmallName}} from storage based on criteria
func (u *{{.UcFirstName}}Usecase) Fetch(ctx context.Context, ctr *domain.{{.UcFirstName}}Criteria) ([]*domain.{{.UcFirstName}}, error) {
	return u.{{.SmallName}}Repository.Fetch(ctx, ctr)
}

// Count return count of {{.SmallName}} from storage based on criteria
func (u *{{.UcFirstName}}Usecase) Count(ctx context.Context, ctr *domain.{{.UcFirstName}}Criteria) (int64, error) {
	return u.{{.SmallName}}Repository.Count(ctx, ctr)
}

// FetchOne fetch a {{.SmallName}} by primary id
func (u *{{.UcFirstName}}Usecase) FetchOne(ctx context.Context, ctr *domain.{{.UcFirstName}}Criteria) (*domain.{{.UcFirstName}}, error) {
	return u.{{.SmallName}}Repository.FetchOne(ctx, ctr)
}

// Update update a {{.SmallName}} record
func (u *{{.UcFirstName}}Usecase) Update(ctx context.Context, {{.SmallName}} *domain.{{.UcFirstName}}) error {
	{{.SmallName}}.PopulateUpdateTimeStamp() // update the updated_at timestamp
	return u.{{.SmallName}}Repository.Update(ctx, {{.SmallName}})
}

// Delete soft delete a {{.SmallName}} record
func (u *{{.UcFirstName}}Usecase) Delete(ctx context.Context, ctr *domain.{{.UcFirstName}}Criteria) error {
	return u.{{.SmallName}}Repository.Delete(ctx, ctr)
}
`
