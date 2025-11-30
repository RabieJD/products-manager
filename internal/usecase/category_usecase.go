package usecase

import (
	"strings"

	"github.com/mytheresa/go-hiring-challenge/domain/entity"
	"github.com/mytheresa/go-hiring-challenge/domain/repository"
)

// CategoryUsecase handles category related business logic.
type CategoryUsecase struct {
	categoryRepo repository.CategoryRepository
}

// NewCategoryUsecase creates a new category usecase.
func NewCategoryUsecase(categoryRepo repository.CategoryRepository) *CategoryUsecase {
	return &CategoryUsecase{
		categoryRepo: categoryRepo,
	}
}

// GetAll retrieves every category.
func (u *CategoryUsecase) GetAll() ([]entity.Category, error) {
	return u.categoryRepo.GetAllCategories()
}

// Create creates a category after trimming inputs.
func (u *CategoryUsecase) Create(code, name string) (*entity.Category, error) {
	category := &entity.Category{
		Code: strings.TrimSpace(code),
		Name: strings.TrimSpace(name),
	}
	return u.categoryRepo.CreateCategory(category)
}
