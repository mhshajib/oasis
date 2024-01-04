package cli_template

var Domain string = `
package domain

import (
    "context"
    "errors"
)

// {{.UcFirstName}} represents a {{.SmallName}}
type {{.UcFirstName}} struct {
    ID          string              ` + "`json:\"_id,omitempty\"` " + `
    FieldOne    string             	` + "`json:\"field_one\"`" + `
    TimeStamp
}

// {{.UcFirstName}}Criteria represents criteria for filtering out {{.SmallName}}
type {{.UcFirstName}}Criteria struct {
    ID            *string
    FieldOne      *string
    Offset, Limit *int64
    WithDeleted   *bool
    SortAsc       bool
}

// {{.UcFirstName}}Repository represents {{.SmallName}}'s repository contract
type {{.UcFirstName}}Repository interface {
    Store(ctx context.Context, {{.SmallName}} *{{.UcFirstName}}) (*{{.UcFirstName}}, error)
    Fetch(ctx context.Context, ctr *{{.UcFirstName}}Criteria) ([]*{{.UcFirstName}}, error)
    FetchOne(ctx context.Context, ctr *{{.UcFirstName}}Criteria) (*{{.UcFirstName}}, error)
    Count(ctx context.Context, ctr *{{.UcFirstName}}Criteria) (int64, error)
    Update(ctx context.Context, {{.SmallName}} *{{.UcFirstName}}) (*{{.UcFirstName}}, error)
    Delete(ctx context.Context, ctr *{{.UcFirstName}}Criteria) error
}

// {{.UcFirstName}}Usecase represents {{.SmallName}}'s usecase contract
type {{.UcFirstName}}Usecase interface {
    Store(ctx context.Context, {{.SmallName}} *{{.UcFirstName}}) (*{{.UcFirstName}}, error)
    Fetch(ctx context.Context, ctr *{{.UcFirstName}}Criteria) ([]*{{.UcFirstName}}, error)
    FetchOne(ctx context.Context, ctr *{{.UcFirstName}}Criteria) (*{{.UcFirstName}}, error)
    Update(ctx context.Context, {{.SmallName}} *{{.UcFirstName}}) (*{{.UcFirstName}}, error)
    Delete(ctx context.Context, ctr *{{.UcFirstName}}Criteria) error
}

// common available {{.SmallName}} errors
var (
    Err{{.UcFirstName}}NotFound = errors.New("{{.SmallName}} not found")
    Err{{.UcFirstName}}Conflict = errors.New("{{.SmallName}} conflict")
)
`
