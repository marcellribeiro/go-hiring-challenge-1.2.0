package product

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/internal/repository"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProductsRepository is a mock implementation of ProductsInterface
type MockProductsRepository struct {
	mock.Mock
}

func (m *MockProductsRepository) GetProducts(filter repository.ProductsFilter) ([]models.Product, int64, error) {
	args := m.Called(filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductsRepository) GetProductByCode(code string) (*models.Product, error) {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func TestHandleGetByCode(t *testing.T) {
	t.Run("returns product with category and variants", func(t *testing.T) {
		mockRepo := new(MockProductsRepository)
		handler := NewProductHandler(mockRepo)

		categoryID := uint(1)
		product := &models.Product{
			ID:         1,
			Code:       "PROD001",
			Price:      decimal.NewFromFloat(10.99),
			CategoryID: &categoryID,
			Category: &models.Category{
				Code: "CLOTHING",
				Name: "Clothing",
			},
			Variants: []models.Variant{
				{
					Name:  "Variant A",
					SKU:   "SKU001A",
					Price: decimal.NewFromFloat(11.99),
				},
				{
					Name:  "Variant B",
					SKU:   "SKU001B",
					Price: decimal.Zero,
				},
			},
		}

		mockRepo.On("GetProductByCode", "PROD001").Return(product, nil)

		req := httptest.NewRequest(http.MethodGet, "/catalog/PROD001", nil)
		req.SetPathValue("code", "PROD001")
		rec := httptest.NewRecorder()

		handler.HandleGetByCode(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response models.Product
		json.NewDecoder(rec.Body).Decode(&response)
		assert.Equal(t, "PROD001", response.Code)
		assert.NotNil(t, response.Category)
		assert.Equal(t, "CLOTHING", response.Category.Code)
		assert.Len(t, response.Variants, 2)

		mockRepo.AssertExpectations(t)
	})

	t.Run("returns 404 when product not found", func(t *testing.T) {
		mockRepo := new(MockProductsRepository)
		handler := NewProductHandler(mockRepo)

		mockRepo.On("GetProductByCode", "NOTFOUND").Return(nil, errors.New("not found"))

		req := httptest.NewRequest(http.MethodGet, "/catalog/NOTFOUND", nil)
		req.SetPathValue("code", "NOTFOUND")
		rec := httptest.NewRecorder()

		handler.HandleGetByCode(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		mockRepo.AssertExpectations(t)
	})
}
