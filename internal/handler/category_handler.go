package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/internal/dto"
	"github.com/mytheresa/go-hiring-challenge/internal/usecase"
)

// CategoryHandler manages category endpoints.
type CategoryHandler struct {
	categoryUsecase *usecase.CategoryUsecase
}

// NewCategoryHandler constructs a CategoryHandler.
func NewCategoryHandler(categoryUsecase *usecase.CategoryUsecase) *CategoryHandler {
	return &CategoryHandler{categoryUsecase: categoryUsecase}
}

// HandleGetCategories handles GET /categories.
func (h *CategoryHandler) HandleGetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.categoryUsecase.GetAll()
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := make([]dto.CategoryResponse, len(categories))
	for i, c := range categories {
		response[i] = dto.CategoryResponse{
			Code: c.Code,
			Name: c.Name,
		}
	}

	api.OKResponse(w, response)
}

// HandleCreateCategory handles POST /categories.
func (h *CategoryHandler) HandleCreateCategory(w http.ResponseWriter, r *http.Request) {
	var payload dto.CategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	payload.Code = strings.TrimSpace(payload.Code)
	payload.Name = strings.TrimSpace(payload.Name)

	if payload.Code == "" || payload.Name == "" {
		api.ErrorResponse(w, http.StatusBadRequest, "code and name are required")
		return
	}

	category, err := h.categoryUsecase.Create(payload.Code, payload.Name)
	if err != nil {
		h.handleCreateError(w, err)
		return
	}

	response := dto.CategoryResponse{
		Code: category.Code,
		Name: category.Name,
	}

	api.CreatedResponse(w, response)
}

func (h *CategoryHandler) handleCreateError(w http.ResponseWriter, err error) {
	var httpStatus = http.StatusInternalServerError

	// Placeholder for future error mapping (e.g., unique violations)
	var validationErr interface{ Error() string }
	if errors.As(err, &validationErr) {
		// currently falls back to 500
	}

	api.ErrorResponse(w, httpStatus, err.Error())
}
