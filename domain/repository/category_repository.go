package repository

import "github.com/mytheresa/go-hiring-challenge/domain/entity"

// CategoryRepository defines operations related to categories.
type CategoryRepository interface {
	// GetAllCategories retrieves all categories.
	GetAllCategories() ([]entity.Category, error)

	// CreateCategory persists a new category.
	CreateCategory(category *entity.Category) (*entity.Category, error)
}
