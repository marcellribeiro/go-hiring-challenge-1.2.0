package product

import (
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/internal/repository"
)

type ProductHandler struct {
	repo repository.ProductsInterface
}

func NewProductHandler(r repository.ProductsInterface) *ProductHandler {
	return &ProductHandler{
		repo: r,
	}
}

func (h *ProductHandler) HandleGetByCode(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		api.ErrorResponse(w, http.StatusBadRequest, "Product code is required")
		return
	}

	// Get product from repository
	product, err := h.repo.GetProductByCode(code)
	if err != nil {
		api.ErrorResponse(w, http.StatusNotFound, "Product not found")
		return
	}

	// Return the product as a JSON response
	api.OKResponse(w, product)
}
