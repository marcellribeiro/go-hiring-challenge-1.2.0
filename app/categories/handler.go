package categories

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/internal/common"
	"github.com/mytheresa/go-hiring-challenge/internal/repository"
	"github.com/mytheresa/go-hiring-challenge/models"
)

type CategoriesResponse struct {
	Categories []models.Category `json:"categories"`
	Total      int64             `json:"total"`
}

type CreateCategoryRequest struct {
	Code string `json:"code" validate:"required"`
	Name string `json:"name" validate:"required"`
}

type CategoriesHandler struct {
	repo repository.CategoriesInterface
}

func NewCategoriesHandler(r repository.CategoriesInterface) *CategoriesHandler {
	return &CategoriesHandler{
		repo: r,
	}
}

func (h *CategoriesHandler) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	// Process filters from request
	filter := h.processFilters(r)

	// Get categories from repository
	categories, total, err := h.repo.GetAllCategories(filter)
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Return the categories as a JSON response
	response := CategoriesResponse{
		Categories: categories,
		Total:      total,
	}
	api.OKResponse(w, response)
}

func (h *CategoriesHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var req CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	var validate = validator.New()
	if err := validate.Struct(req); err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "Validation error: "+err.Error())
		return
	}

	// Create category
	category := &models.Category{
		Code: req.Code,
		Name: req.Name,
	}

	if err := h.repo.CreateCategory(category); err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Return created category
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

func (h *CategoriesHandler) processFilters(r *http.Request) repository.CategoriesFilter {
	// Parse pagination parameters
	offset, limit := common.ParseOffsetLimit(r)

	// Build filter
	filter := repository.CategoriesFilter{
		Offset: offset,
		Limit:  limit,
	}

	return filter
}
