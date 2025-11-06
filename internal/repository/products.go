package repository

import (
	"github.com/mytheresa/go-hiring-challenge/internal/database"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
)

type ProductsInterface interface {
	GetProducts(filter ProductsFilter) ([]models.Product, int64, error)
	GetProductByCode(code string) (*models.Product, error)
}

type Products struct {
	db database.Database
}

type ProductsFilter struct {
	CategoryCode *string
	MaxPrice     *decimal.Decimal
	Offset       int
	Limit        int
}

func NewProducts(db database.Database) *Products {
	return &Products{
		db: db,
	}
}

func (r *Products) GetProducts(filter ProductsFilter) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := r.db.Model(&models.Product{})

	// Apply category filter
	if filter.CategoryCode != nil {
		query = query.Joins("JOIN categories ON categories.id = products.category_id").
			Where("categories.code = ?", *filter.CategoryCode)
	}

	// Apply max price filter
	if filter.MaxPrice != nil {
		query = query.Where("products.price <= ?", *filter.MaxPrice)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and preload relations
	if err := query.
		Preload("Category").
		Preload("Variants").
		Offset(filter.Offset).
		Limit(filter.Limit).
		Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
func (r *Products) GetProductByCode(code string) (*models.Product, error) {
	var product models.Product
	if err := r.db.Where("code = ?", code).
		Preload("Category").
		Preload("Variants").
		First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}
