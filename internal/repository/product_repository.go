package repository

import (
	"errors"

	"github.com/mytheresa/go-hiring-challenge/domain/entity"
	domainRepo "github.com/mytheresa/go-hiring-challenge/domain/repository"
	"github.com/mytheresa/go-hiring-challenge/internal/repository/model"
	"gorm.io/gorm"
)

// productRepository implements the ProductRepository interface.
type productRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new product repository instance.
func NewProductRepository(db *gorm.DB) domainRepo.ProductRepository {
	return &productRepository{
		db: db,
	}
}

// Ensure productRepository implements ProductRepository interface at compile time.
var _ domainRepo.ProductRepository = (*productRepository)(nil)

// GetProducts retrieves products with pagination and filtering.
func (r *productRepository) GetProducts(filter *domainRepo.ProductFilter, pagination *domainRepo.PaginationParams) (*domainRepo.PaginationResult, error) {
	// Build base query for counting (without preloads for performance)
	countQuery := r.db.Model(&model.ProductDB{})

	// Build query for fetching (with preloads)
	query := r.db.Model(&model.ProductDB{}).Preload("Variants").Preload("Category")

	// Apply filters
	if filter != nil {
		if filter.CategoryCode != nil && *filter.CategoryCode != "" {
			// Use subquery to filter by category code
			countQuery = countQuery.Joins("JOIN product_categories ON products.category_id = product_categories.id").
				Where("product_categories.code = ?", *filter.CategoryCode)
			query = query.Joins("JOIN product_categories ON products.category_id = product_categories.id").
				Where("product_categories.code = ?", *filter.CategoryCode)
		}
		if filter.PriceLessThan != nil {
			countQuery = countQuery.Where("price < ?", *filter.PriceLessThan)
			query = query.Where("price < ?", *filter.PriceLessThan)
		}
	}

	// Get total count
	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, err
	}

	// Apply pagination
	if pagination != nil {
		query = query.Offset(pagination.Offset).Limit(pagination.Limit)
	}

	// Fetch products
	var productsDB []model.ProductDB
	if err := query.Find(&productsDB).Error; err != nil {
		return nil, err
	}

	products := make([]entity.Product, len(productsDB))
	for i, p := range productsDB {
		products[i] = *r.toEntity(&p)
	}

	return &domainRepo.PaginationResult{
		Products: products,
		Total:    total,
	}, nil
}

// GetProductByCode retrieves a product by its code.
func (r *productRepository) GetProductByCode(code string) (*entity.Product, error) {
	var productDB model.ProductDB
	if err := r.db.Preload("Variants").Preload("Category").Where("code = ?", code).First(&productDB).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return r.toEntity(&productDB), nil
}

// GetProductByID retrieves a product by its ID.
func (r *productRepository) GetProductByID(id uint) (*entity.Product, error) {
	var productDB model.ProductDB
	if err := r.db.Preload("Variants").Preload("Category").First(&productDB, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return r.toEntity(&productDB), nil
}

// toEntity converts a database model to a domain entity.
func (r *productRepository) toEntity(p *model.ProductDB) *entity.Product {
	if p == nil {
		return nil
	}

	variants := make([]entity.Variant, len(p.Variants))
	for i, v := range p.Variants {
		variants[i] = entity.Variant{
			ID:        v.ID,
			ProductID: v.ProductID,
			Name:      v.Name,
			SKU:       v.SKU,
			Price:     v.Price,
		}
	}

	var category *entity.Category
	if p.Category != nil {
		category = &entity.Category{
			ID:   p.Category.ID,
			Code: p.Category.Code,
			Name: p.Category.Name,
		}
	}

	return &entity.Product{
		ID:       p.ID,
		Code:     p.Code,
		Price:    p.Price,
		Category: category,
		Variants: variants,
	}
}
