package catalog

import (
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/internal/common"
	"github.com/mytheresa/go-hiring-challenge/internal/repository"
	"github.com/shopspring/decimal"
)

type Response struct {
	Products []interface{} `json:"products"`
	Total    int64         `json:"total"`
}

type CatalogHandler struct {
	repo repository.ProductsInterface
}

func NewCatalogHandler(r repository.ProductsInterface) *CatalogHandler {
	return &CatalogHandler{
		repo: r,
	}
}

func (h *CatalogHandler) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	// Process filters from request
	filter := h.processFilters(r)

	// Get products from repository
	products, total, err := h.repo.GetProducts(filter)
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Return the products as a JSON response
	response := Response{
		Products: make([]interface{}, len(products)),
		Total:    total,
	}

	for i := range products {
		response.Products[i] = products[i]
	}

	api.OKResponse(w, response)
}

func (h *CatalogHandler) processFilters(r *http.Request) repository.ProductsFilter {
	// Parse pagination parameters
	offset, limit := common.ParseOffsetLimit(r)

	// Build filter
	filter := repository.ProductsFilter{
		Offset: offset,
		Limit:  limit,
	}

	// Parse category filter
	if categoryCode := r.URL.Query().Get("category"); categoryCode != "" {
		filter.CategoryCode = &categoryCode
	}

	// Parse max price filter
	if maxPriceStr := r.URL.Query().Get("priceLessThan"); maxPriceStr != "" {
		if maxPrice, err := decimal.NewFromString(maxPriceStr); err == nil {
			filter.MaxPrice = &maxPrice
		}
	}

	return filter
}
