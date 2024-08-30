package cli_template

var HttpHandler string = `
package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"{{.DomainPath}}"
	"{{.ModuleName}}/pkg/log"
	//"{{.ModuleName}}/pkg/middleware"
	"{{.ModuleName}}/pkg/utils"
	"{{.TransformerPath}}"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// {{.UcFirstName}}Handler represents {{.SmallName}} HTTP/JSON handler
type {{.UcFirstName}}Handler struct {
	{{.UcFirstName}}Usecase domain.{{.UcFirstName}}Usecase
}

// New{{.UcFirstName}}Handler will initialize the resources endpoint
func New{{.UcFirstName}}Handler(r *chi.Mux, uc domain.{{.UcFirstName}}Usecase) {
	handler := &{{.UcFirstName}}Handler{
		{{.UcFirstName}}Usecase: uc,
	}
	r.Route("/v1/{{.SmallPluralName}}", func(r chi.Router) {
		// all the routes are protected
		// r.Use(middleware.Auth)

		r.Post("/", handler.Store)
		r.Get("/", handler.Fetch)
		r.Get("/{id}", handler.FetchByID)
		r.Put("/{id}", handler.Update)
		r.Delete("/{id}", handler.Delete)
	})
}

// ReqCreate{{.UcFirstName}} represents create {{.SmallName}} request
type ReqCreate{{.UcFirstName}} struct {  {{range .Fields}}
    {{.Name}}    {{.Type}}    ` + "`json:\"{{.JsonTag}}\"`" + ` {{end}}
}

// Validate validate create {{.SmallName}} requests
func (r *ReqCreate{{.UcFirstName}}) Validate(ctx context.Context) utils.Errors {
	errs := utils.Errors{}

	{{range .Fields}}
		{{if eq .Type "string"}}
			r.{{.Name}} = strings.TrimSpace(r.{{.Name}})

			if r.{{.Name}} == "" {
				errs.Add("{{.JsonTag}}", "is required")
			}
		{{else if or (eq .Type "int") (eq .Type "int32") (eq .Type "int64") (eq .Type "float64") (eq .Type "float32")}}
			if r.{{.Name}} {{.OperatorLessThan}} 0 {
				errs.Add("{{.JsonTag}}", "must be non-negative")
			}
		{{else if or (eq .Type "[]string") (eq .Type "[]int")}}
			if len(r.{{.Name}}) == 0 {
				errs.Add("{{.JsonTag}}", "cannot be empty")
			}
		{{end}}
	{{end}}


	return errs
}

// Store validate {{.SmallName}} input, authorization and store a {{.SmallName}}
func (h *{{.UcFirstName}}Handler) Store(w http.ResponseWriter, r *http.Request) {
    req := ReqCreate{{.UcFirstName}}{}
	ctx := r.Context()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		(&utils.Response{
			Status:  http.StatusBadRequest,
			Message: "Bad request",
			Error:   err.Error(),
		}).Render(w)
		return
	}

	vErr := req.Validate(ctx)
	if !vErr.IsNil() {
		(&utils.Response{
			Status:  http.StatusUnprocessableEntity,
			Message: "Validation error",
			Error:   vErr,
		}).Render(w)
		return
	}

	{{.SmallName}} := domain.{{.UcFirstName}}{ {{range .Fields}}
    	{{.Name}}: req.{{.Name}}, {{end}}
	}

	{{.SmallName}}Resp, err := h.{{.UcFirstName}}Usecase.Store(ctx, &{{.SmallName}})
	if err != nil {
		log.ErrorWithFields("Failed to create {{.SmallName}}", log.Fields{
			"event": "store_{{.SmallName}}_log",
			"error": err.Error(),
		})

		(&utils.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to create {{.SmallName}}",
			Error:   err.Error(),
		}).Render(w)
		return
	}

	(&utils.Response{
		Status:  http.StatusCreated,
		Message: "{{.UcFirstName}} created successfully",
		Data:    transformer.Transform{{.UcFirstName}}({{.SmallName}}Resp),
	}).Render(w)
}

// Fetch return a list of {{.SmallName}} based on criteria
func (h *{{.UcFirstName}}Handler) Fetch(w http.ResponseWriter, r *http.Request) {
	page, limit, offset, err := utils.GetPager(r)
	if err != nil {
		(&utils.Response{
			Status:  http.StatusBadRequest,
			Message: "Failed to fetch {{.SmallName}}",
			Error:   err.Error(),
		}).Render(w)
		return
	}

	ctx := r.Context()

	ctr := &domain.{{.UcFirstName}}Criteria{}
	ctr.Limit = &limit
	ctr.Offset = &offset

	if v := r.URL.Query().Get("id"); v != "" {
		ctr.ID = &v
	}

	{{range .CriteriaFields}}
		{{if eq .Type "*string"}}
			if v := r.URL.Query().Get("{{.JsonTag}}"); v != "" {
				ctr.{{.Name}} = &v
			}
		{{else if or (eq .Type "int") (eq .Type "int32") (eq .Type "int64") (eq .Type "float64") (eq .Type "float32")}}
			if v := r.URL.Query().Get("{{.JsonTag}}"); v != "" {
				val, _ := strconv.Atoi(v)
				ctr.{{.Name}} = val
			}
		{{end}}
	{{end}}

	if v := r.URL.Query().Get("asc"); v != "" {
		val, _ := strconv.ParseBool(v)
		ctr.SortAsc = val
	}

	{{.SmallName}}, err := h.{{.UcFirstName}}Usecase.Fetch(ctx, ctr)
	if err != nil {
		log.ErrorWithFields("Failed to fetch {{.SmallName}}", log.Fields{
			"event": "fetch_{{.SmallName}}",
			"error": err.Error(),
		})

		(&utils.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to fetch {{.SmallName}}",
			Error:   err.Error(),
		}).Render(w)
		return
	}

	pagination := utils.NewPagination(int64(len({{.SmallName}})), page, limit)
	(&utils.Response{
		Status:     http.StatusOK,
		Data:       transformer.Transform{{.UcFirstName}}List({{.SmallName}}),
		Pagination: pagination,
	}).Render(w)
}

// FetchByID return {{.SmallName}} by id
func (h *{{.UcFirstName}}Handler) FetchByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		(&utils.Response{
			Status:  http.StatusBadRequest,
			Message: "Failed to fetch {{.SmallName}}",
			Error:   err.Error(),
		}).Render(w)
		return
	}

	ctx := r.Context()
	ctr := &domain.{{.UcFirstName}}Criteria{
		ID: &id,
	}

	{{.SmallName}}, err := h.{{.UcFirstName}}Usecase.FetchOne(ctx, ctr)
	if err != nil {
		if err == domain.Err{{.UcFirstName}}NotFound {
			(&utils.Response{
				Status:  http.StatusNotFound,
				Message: "{{.UcFirstName}} not found",
				Error:   err.Error(),
			}).Render(w)
			return
		}

		log.ErrorWithFields("Failed to fetch {{.SmallName}}", log.Fields{
			"event": "fetch_{{.SmallName}}",
			"error": err.Error(),
		})

		(&utils.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to fetch {{.SmallName}}",
			Error:   err.Error(),
		}).Render(w)
		return
	}

	(&utils.Response{
		Status: http.StatusOK,
		Data:   transformer.Transform{{.UcFirstName}}({{.SmallName}}),
	}).Render(w)
}
// ReqUpdate{{.UcFirstName}} represents create {{.SmallName}} request
type ReqUpdate{{.UcFirstName}} struct { {{range .Fields}}
    {{.Name}}    {{.Type}}    ` + "`json:\"{{.JsonTag}}\"`" + ` {{end}}
}

// Validate validate create {{.SmallName}} requests
func (r *ReqUpdate{{.UcFirstName}}) Validate(ctx context.Context) utils.Errors {
	errs := utils.Errors{}

	{{range .Fields}}
		{{if eq .Type "string"}}
			r.{{.Name}} = strings.TrimSpace(r.{{.Name}})

			if r.{{.Name}} == "" {
				errs.Add("{{.JsonTag}}", "is required")
			}
		{{else if or (eq .Type "int") (eq .Type "int32") (eq .Type "int64") (eq .Type "float64") (eq .Type "float32")}}
			if r.{{.Name}} {{.OperatorLessThan}} 0 {
				errs.Add("{{.JsonTag}}", "must be non-negative")
			}
		{{else if or (eq .Type "[]string") (eq .Type "[]int")}}
			if len(r.{{.Name}}) == 0 {
				errs.Add("{{.JsonTag}}", "cannot be empty")
			}
		{{end}}
	{{end}}

	return errs
}

// Update modify {{.SmallName}} record
func (h *{{.UcFirstName}}Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		(&utils.Response{
			Status:  http.StatusBadRequest,
			Message: "Failed to fetch {{.SmallName}}",
			Error:   err.Error(),
		}).Render(w)
		return
	}

	req := ReqUpdate{{.UcFirstName}}{}
	ctx := r.Context()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		(&utils.Response{
			Status:  http.StatusBadRequest,
			Message: "Bad request",
			Error:   err.Error(),
		}).Render(w)
		return
	}

	vErr := req.Validate(ctx)
	if !vErr.IsNil() {
		(&utils.Response{
			Status:  http.StatusUnprocessableEntity,
			Message: "Validation error",
			Error:   vErr,
		}).Render(w)
		return
	}

	ctr := &domain.{{.UcFirstName}}Criteria{
		ID: &id,
	}

	{{.SmallName}}, err := h.{{.UcFirstName}}Usecase.FetchOne(ctx, ctr)
	if err != nil {
		if err == domain.Err{{.UcFirstName}}NotFound {
			(&utils.Response{
				Status:  http.StatusNotFound,
				Message: "{{.UcFirstName}} not found",
				Error:   err.Error(),
			}).Render(w)
			return
		}

		log.ErrorWithFields("Failed to fetch {{.SmallName}}", log.Fields{
			"event": "update_{{.SmallName}}",
			"error": err.Error(),
		})

		(&utils.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to fetch {{.SmallName}}",
			Error:   err.Error(),
		}).Render(w)
		return
	}

	{{range .Fields}}
		{{if eq .Type "string"}}
			if v := req.{{.Name}}; v != "" {
				{{$.SmallName}}.{{.Name}} = v
			}
		{{else if or (eq .Type "int") (eq .Type "int32") (eq .Type "int64") (eq .Type "float64") (eq .Type "float32")}}
			if v := req.{{.Name}}; v {{.OperatorGreterThan}} 0 {
				{{$.SmallName}}.{{.Name}} = v
			}
		{{else if or (eq .Type "[]string") (eq .Type "[]int")}}
			if len(r.{{.Name}}) == 0 {
				errs.Add("{{.JsonTag}}", "cannot be empty")
			}
		{{end}}
	{{end}}

	{{.SmallName}}Resp, err := h.{{.UcFirstName}}Usecase.Update(ctx, {{.SmallName}})
	if err != nil {
		log.ErrorWithFields("Failed to update {{.SmallName}}", log.Fields{
			"event": "update_{{.SmallName}}",
			"error": err.Error(),
		})

		(&utils.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to update {{.SmallName}}",
			Error:   err.Error(),
		}).Render(w)
		return
	}

	(&utils.Response{
		Status: http.StatusOK,
		Data:   transformer.Transform{{.UcFirstName}}({{.SmallName}}Resp),
	}).Render(w)
}

// Delete {{.SmallName}} record
func (h {{.UcFirstName}}Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		(&utils.Response{
			Status:  http.StatusBadRequest,
			Message: "Failed to fetch {{.SmallName}}",
			Error:   err.Error(),
		}).Render(w)
		return
	}

	err = h.{{.UcFirstName}}Usecase.Delete(r.Context(), &domain.{{.UcFirstName}}Criteria{
		ID: &id,
	})
	if err != nil && err != domain.Err{{.UcFirstName}}NotFound {
		log.ErrorWithFields("Failed to delete {{.SmallName}}", log.Fields{
			"event": "delete_{{.SmallName}}",
			"error": err.Error(),
		})

		_ = (&utils.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to delete {{.SmallName}}",
			Error:   err.Error(),
		}).Render(w)
		return
	}

	_ = (&utils.Response{
		Status:  http.StatusOK,
		Message: "{{.UcFirstName}} removed successfully",
	}).Render(w)
}
`
