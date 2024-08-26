
package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"oasis/domain"
	"oasis/pkg/log"
	"oasis/pkg/middleware"
	"oasis/pkg/utils"
	"oasis/transformer"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Handler represents  HTTP/JSON handler
type Handler struct {
	Usecase domain.Usecase
}

// NewHandler will initialize the resources endpoint
func NewHandler(r *chi.Mux, uc domain.Usecase) {
	handler := &Handler{
		Usecase: uc,
	}
	r.Route("/v1/", func(r chi.Router) {
		// all the routes are protected
		// r.Use(middleware.Auth)

		r.Post("/", handler.Store)
		r.Get("/", handler.Fetch)
		r.Get("/{id}", handler.FetchByID)
		r.Put("/{id}", handler.Update)
		r.Delete("/{id}", handler.Delete)
	})
}

// ReqCreate represents create  request
type ReqCreate struct {
	FieldOne string `json:"field_one,omitempty"` 
}

// Validate validate create  requests
func (r *ReqCreate) Validate(ctx context.Context) utils.Errors {
	r.FieldOne = strings.TrimSpace(r.FieldOne)

	errs := utils.Errors{}
	if r.FieldOne == "" {
		errs.Add("field_one", "is required")
	}

	return errs
}

// Store validate  input, authorization and store a 
func (h *Handler) Store(w http.ResponseWriter, r *http.Request) {
    req := ReqCreate{}
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

	 := domain.{
		FieldOne: req.FieldOne,
	}

	Resp, err := h.Usecase.Store(ctx, &)
	if err != nil {
		log.ErrorWithFields("Failed to create ", log.Fields{
			"event": "store__log",
			"error": err.Error(),
		})

		(&utils.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to create ",
			Error:   err.Error(),
		}).Render(w)
		return
	}

	(&utils.Response{
		Status:  http.StatusCreated,
		Message: " created successfully",
		Data:    transformer.TransformClient(Resp),
	}).Render(w)
}

// Fetch return a list of  based on criteria
func (h *Handler) Fetch(w http.ResponseWriter, r *http.Request) {
	page, limit, offset, err := utils.GetPager(r)
	if err != nil {
		(&utils.Response{
			Status:  http.StatusBadRequest,
			Message: "Failed to fetch ",
			Error:   err.Error(),
		}).Render(w)
		return
	}

	ctx := r.Context()

	ctr := &domain.Criteria{}
	ctr.Limit = &limit
	ctr.Offset = &offset

	if v := r.URL.Query().Get("id"); v != "" {
		ctr.ID = &v
	}

	if v := r.URL.Query().Get("field_one"); v != "" {
		ctr.FieldOne = &v
	}

	if v := r.URL.Query().Get("asc"); v != "" {
		val, _ := strconv.ParseBool(v)
		ctr.SortAsc = val
	}

	, err := h.Usecase.Fetch(ctx, ctr)
	if err != nil {
		log.ErrorWithFields("Failed to fetch ", log.Fields{
			"event": "fetch_",
			"error": err.Error(),
		})

		(&utils.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to fetch ",
			Error:   err.Error(),
		}).Render(w)
		return
	}

	pagination := utils.NewPagination(int64(len()), page, limit)
	(&utils.Response{
		Status:     http.StatusOK,
		Data:       transformer.TransformList(),
		Pagination: pagination,
	}).Render(w)
}

// FetchByID return  by id
func (h *Handler) FetchByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		(&utils.Response{
			Status:  http.StatusBadRequest,
			Message: "Failed to fetch ",
			Error:   err.Error(),
		}).Render(w)
		return
	}

	ctx := r.Context()
	ctr := &domain.Criteria{
		ID: &id,
	}

	, err := h.Usecase.FetchOne(ctx, ctr)
	if err != nil {
		if err == domain.ErrNotFound {
			(&utils.Response{
				Status:  http.StatusNotFound,
				Message: " not found",
				Error:   err.Error(),
			}).Render(w)
			return
		}

		log.ErrorWithFields("Failed to fetch ", log.Fields{
			"event": "fetch_",
			"error": err.Error(),
		})

		(&utils.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to fetch ",
			Error:   err.Error(),
		}).Render(w)
		return
	}

	(&utils.Response{
		Status: http.StatusOK,
		Data:   transformer.Transform(),
	}).Render(w)
}
// ReqUpdate represents create  request
type ReqUpdate struct {
	FieldOne string `json:"field_one,omitempty"` 
}

// Validate validate create  requests
func (r *ReqUpdate) Validate(ctx context.Context) utils.Errors {
	r.FieldOne = strings.TrimSpace(r.FieldOne)

	errs := utils.Errors{}
	if r.FieldOne == "" {
		errs.Add("field_one", "is required")
	}

	return errs
}

// Update modify  record
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		(&utils.Response{
			Status:  http.StatusBadRequest,
			Message: "Failed to fetch ",
			Error:   err.Error(),
		}).Render(w)
		return
	}

	req := ReqUpdate{}
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

	ctr := &domain.Criteria{
		ID: &id,
	}

	, err := h.Usecase.FetchOne(ctx, ctr)
	if err != nil {
		if err == domain.ErrNotFound {
			(&utils.Response{
				Status:  http.StatusNotFound,
				Message: " not found",
				Error:   err.Error(),
			}).Render(w)
			return
		}

		log.ErrorWithFields("Failed to fetch ", log.Fields{
			"event": "update_",
			"error": err.Error(),
		})

		(&utils.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to fetch ",
			Error:   err.Error(),
		}).Render(w)
		return
	}

	if v := req.FieldOne; v != "" {
		.FieldOne = v
	}

	Resp, err := h.Usecase.Update(ctx, )
	if err != nil {
		log.ErrorWithFields("Failed to update ", log.Fields{
			"event": "update_",
			"error": err.Error(),
		})

		(&utils.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to update ",
			Error:   err.Error(),
		}).Render(w)
		return
	}

	(&utils.Response{
		Status: http.StatusOK,
		Data:   transformer.Transform(Resp),
	}).Render(w)
}

// Delete  record
func (h Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		(&utils.Response{
			Status:  http.StatusBadRequest,
			Message: "Failed to fetch ",
			Error:   err.Error(),
		}).Render(w)
		return
	}

	err = h.Usecase.Delete(r.Context(), &domain.Criteria{
		ID: &id,
	})
	if err != nil && err != domain.ErrNotFound {
		log.ErrorWithFields("Failed to delete ", log.Fields{
			"event": "delete_",
			"error": err.Error(),
		})

		_ = (&utils.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to delete ",
			Error:   err.Error(),
		}).Render(w)
		return
	}

	_ = (&utils.Response{
		Status:  http.StatusOK,
		Message: " removed successfully",
	}).Render(w)
}
