package repository

import (
	"github.com/mytheresa/go-hiring-challenge/internal/database"
	"github.com/mytheresa/go-hiring-challenge/models"
)

type CategoriesInterface interface {
	GetAllCategories(filter CategoriesFilter) ([]models.Category, int64, error)
	CreateCategory(category *models.Category) error
}

type Categories struct {
	db database.Database
}

type CategoriesFilter struct {
	Offset int
	Limit  int
}

func NewCategories(db database.Database) *Categories {
	return &Categories{
		db: db,
	}
}

func (r *Categories) GetAllCategories(filter CategoriesFilter) ([]models.Category, int64, error) {
	var categories []models.Category
	var total int64

	query := r.db.Model(&models.Category{})

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if err := query.
		Offset(filter.Offset).
		Limit(filter.Limit).
		Find(&categories).Error; err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

func (r *Categories) CreateCategory(category *models.Category) error {
	return r.db.Create(category).Error
}
