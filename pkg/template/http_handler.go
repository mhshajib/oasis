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
	"{{.ModuleName}}/domain"
	"{{.ModuleName}}/pkg/log"
	"{{.ModuleName}}/pkg/middleware"
	"{{.ModuleName}}/pkg/utils"
	"{{.ModuleName}}/{{.SmallName}}/transformer"
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
		r.Use(middleware.Auth)

		r.Post("/", handler.Store)
		r.Get("/", handler.Fetch)
		r.Get("/{id}", handler.FetchByID)
		r.Put("/{id}", handler.Update)
		r.Put("/{id}", handler.Delete)
	})
}

// ReqCreate{{.UcFirstName}} represents create {{.SmallName}} request
type ReqCreate{{.UcFirstName}} struct {
	Field string `json:"field,omitempty"`
}

// Validate validate create {{.SmallName}} requests
func (r *ReqCreate{{.UcFirstName}}) Validate(ctx context.Context) utils.Errors {
	r.Field = strings.TrimSpace(r.Field)

	errs := utils.Errors{}
	if r.Field == "" {
		errs.Add("field", "is required")
	}

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

	var jsonResponse map[string]interface{}
	err := json.Unmarshal([]byte(req.Response), &jsonResponse)
	if err != nil {
		(&utils.Response{
			Status:  http.StatusUnprocessableEntity,
			Message: "Invalid Response JSON",
			Error:   vErr,
		}).Render(w)
		return
	}

	{{.SmallName}} := domain.{{.UcFirstName}}{
		Field: req.Field,
	}

	if err := h.{{.UcFirstName}}Usecase.Store(ctx, &{{.SmallName}}); err != nil {
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
		Data:    {{.SmallName}}.ID,
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

	if v := r.URL.Query().Get("field"); v != "" {
		ctr.Field = &v
	}

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
type ReqUpdate{{.UcFirstName}} struct {
	Field string `json:"field,omitempty"`
}

// Validate validate create {{.SmallName}} requests
func (r *ReqUpdate{{.UcFirstName}}) Validate(ctx context.Context) utils.Errors {
	r.Field = strings.TrimSpace(r.Field)

	errs := utils.Errors{}
	if r.Field == "" {
		errs.Add("field", "is required")
	}

	return errs
}

// Update modify {{.SmallName}} record
func (h *{{.UcFirstName}}Handler) Update(w http.ResponseWriter, r *http.Request) {
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

	if v := req.Field; v != "" {
		{{.SmallName}}.Field = v
	}

	if err := h.{{.UcFirstName}}Usecase.Update(ctx, {{.SmallName}}); err != nil {
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
		Data:   transformer.Transform{{.UcFirstName}}({{.SmallName}}),
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