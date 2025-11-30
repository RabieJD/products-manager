package usecase

import (
	"github.com/mytheresa/go-hiring-challenge/domain/entity"
	"github.com/mytheresa/go-hiring-challenge/domain/repository"
)

// ProductUsecase handles business logic for products.
type ProductUsecase struct {
	productRepo repository.ProductRepository
}

// NewProductUsecase creates a new product usecase instance.
func NewProductUsecase(productRepo repository.ProductRepository) *ProductUsecase {
	return &ProductUsecase{
		productRepo: productRepo,
	}
}

// GetProducts retrieves products with pagination and filtering.
func (u *ProductUsecase) GetProducts(filter *repository.ProductFilter, pagination *repository.PaginationParams) (*repository.PaginationResult, error) {
	return u.productRepo.GetProducts(filter, pagination)
}

// GetProductByCode retrieves a product by its code.
func (u *ProductUsecase) GetProductByCode(code string) (*entity.Product, error) {
	return u.productRepo.GetProductByCode(code)
}

// GetProductByID retrieves a product by its ID.
func (u *ProductUsecase) GetProductByID(id uint) (*repository.PaginationResult, error) {
	product, err := u.productRepo.GetProductByID(id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return &repository.PaginationResult{
			Products: []entity.Product{},
			Total:    0,
		}, nil
	}
	return &repository.PaginationResult{
		Products: []entity.Product{*product},
		Total:    1,
	}, nil
}
