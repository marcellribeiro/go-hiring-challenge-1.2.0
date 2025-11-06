package catalog

import (
	"encoding/json"
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

func TestHandleGetAll(t *testing.T) {
	t.Run("returns products with category and total count", func(t *testing.T) {
		mockRepo := new(MockProductsRepository)
		handler := NewCatalogHandler(mockRepo)

		categoryID := uint(1)
		products := []models.Product{
			{
				Code:       "PROD001",
				Price:      decimal.NewFromFloat(10.99),
				CategoryID: &categoryID,
				Category: &models.Category{
					Code: "CLOTHING",
					Name: "Clothing",
				},
			},
		}

		mockRepo.On("GetProducts", mock.Anything).Return(products, int64(1), nil)

		req := httptest.NewRequest(http.MethodGet, "/catalog", nil)
		rec := httptest.NewRecorder()

		handler.HandleGetAll(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response Response
		json.NewDecoder(rec.Body).Decode(&response)
		assert.Len(t, response.Products, 1)
		assert.Equal(t, int64(1), response.Total)

		mockRepo.AssertExpectations(t)
	})

	t.Run("filters by category", func(t *testing.T) {
		mockRepo := new(MockProductsRepository)
		handler := NewCatalogHandler(mockRepo)

		mockRepo.On("GetProducts", mock.MatchedBy(func(filter repository.ProductsFilter) bool {
			return filter.CategoryCode != nil && *filter.CategoryCode == "SHOES"
		})).Return([]models.Product{}, int64(0), nil)

		req := httptest.NewRequest(http.MethodGet, "/catalog?category=SHOES", nil)
		rec := httptest.NewRecorder()

		handler.HandleGetAll(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("filters by price less than", func(t *testing.T) {
		mockRepo := new(MockProductsRepository)
		handler := NewCatalogHandler(mockRepo)

		mockRepo.On("GetProducts", mock.MatchedBy(func(filter repository.ProductsFilter) bool {
			return filter.MaxPrice != nil && filter.MaxPrice.Equal(decimal.NewFromFloat(50.00))
		})).Return([]models.Product{}, int64(0), nil)

		req := httptest.NewRequest(http.MethodGet, "/catalog?priceLessThan=50.00", nil)
		rec := httptest.NewRecorder()

		handler.HandleGetAll(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("supports pagination with offset and limit", func(t *testing.T) {
		mockRepo := new(MockProductsRepository)
		handler := NewCatalogHandler(mockRepo)

		mockRepo.On("GetProducts", mock.MatchedBy(func(filter repository.ProductsFilter) bool {
			return filter.Offset == 5 && filter.Limit == 20
		})).Return([]models.Product{}, int64(0), nil)

		req := httptest.NewRequest(http.MethodGet, "/catalog?offset=5&limit=20", nil)
		rec := httptest.NewRecorder()

		handler.HandleGetAll(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockRepo.AssertExpectations(t)
	})
}
