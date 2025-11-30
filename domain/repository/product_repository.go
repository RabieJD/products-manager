package repository

import (
	"github.com/mytheresa/go-hiring-challenge/domain/entity"
	"github.com/shopspring/decimal"
)

// ProductFilter represents filtering criteria for products.
type ProductFilter struct {
	CategoryCode *string
	PriceLessThan *decimal.Decimal
}

// PaginationParams represents pagination parameters.
type PaginationParams struct {
	Offset int
	Limit  int
}

// PaginationResult represents the result of a paginated query.
type PaginationResult struct {
	Products []entity.Product
	Total    int64
}

// ProductRepository defines the interface for product data access operations.
type ProductRepository interface {
	// GetProducts retrieves products with pagination and filtering.
	GetProducts(filter *ProductFilter, pagination *PaginationParams) (*PaginationResult, error)
	
	// GetProductByCode retrieves a product by its code.
	GetProductByCode(code string) (*entity.Product, error)
	
	// GetProductByID retrieves a product by its ID.
	GetProductByID(id uint) (*entity.Product, error)
}

