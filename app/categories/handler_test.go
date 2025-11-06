package categories

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/internal/repository"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCategoriesRepository is a mock implementation of CategoriesInterface
type MockCategoriesRepository struct {
	mock.Mock
}

func (m *MockCategoriesRepository) GetAllCategories(filter repository.CategoriesFilter) ([]models.Category, int64, error) {
	args := m.Called(filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Category), args.Get(1).(int64), args.Error(2)
}

func (m *MockCategoriesRepository) CreateCategory(category *models.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func TestHandleGetAll(t *testing.T) {
	t.Run("returns all categories", func(t *testing.T) {
		mockRepo := new(MockCategoriesRepository)
		handler := NewCategoriesHandler(mockRepo)

		categories := []models.Category{
			{Code: "CLOTHING", Name: "Clothing"},
			{Code: "SHOES", Name: "Shoes"},
			{Code: "ACCESSORIES", Name: "Accessories"},
		}

		mockRepo.On("GetAllCategories", mock.Anything).Return(categories, int64(3), nil)

		req := httptest.NewRequest(http.MethodGet, "/categories", nil)
		rec := httptest.NewRecorder()

		handler.HandleGetAll(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response CategoriesResponse
		json.NewDecoder(rec.Body).Decode(&response)
		assert.Len(t, response.Categories, 3)
		assert.Equal(t, int64(3), response.Total)

		mockRepo.AssertExpectations(t)
	})
}

func TestHandleCreate(t *testing.T) {
	t.Run("creates new category", func(t *testing.T) {
		mockRepo := new(MockCategoriesRepository)
		handler := NewCategoriesHandler(mockRepo)

		reqBody := CreateCategoryRequest{
			Code: "ELECTRONICS",
			Name: "Electronics",
		}

		mockRepo.On("CreateCategory", mock.MatchedBy(func(cat *models.Category) bool {
			return cat.Code == "ELECTRONICS" && cat.Name == "Electronics"
		})).Return(nil)

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		handler.HandleCreate(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var response models.Category
		json.NewDecoder(rec.Body).Decode(&response)
		assert.Equal(t, "ELECTRONICS", response.Code)

		mockRepo.AssertExpectations(t)
	})

	t.Run("returns error for missing required fields", func(t *testing.T) {
		mockRepo := new(MockCategoriesRepository)
		handler := NewCategoriesHandler(mockRepo)

		reqBody := CreateCategoryRequest{
			Code: "ELECTRONICS",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		handler.HandleCreate(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
