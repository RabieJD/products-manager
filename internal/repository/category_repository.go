package repository

import (
	"github.com/mytheresa/go-hiring-challenge/domain/entity"
	domainRepo "github.com/mytheresa/go-hiring-challenge/domain/repository"
	"github.com/mytheresa/go-hiring-challenge/internal/repository/model"
	"gorm.io/gorm"
)

// categoryRepository implements CategoryRepository using GORM.
type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a category repository backed by GORM.
func NewCategoryRepository(db *gorm.DB) domainRepo.CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) GetAllCategories() ([]entity.Category, error) {
	var categoriesDB []model.CategoryDB
	if err := r.db.Find(&categoriesDB).Error; err != nil {
		return nil, err
	}

	categories := make([]entity.Category, len(categoriesDB))
	for i, c := range categoriesDB {
		categories[i] = entity.Category{
			ID:   c.ID,
			Code: c.Code,
			Name: c.Name,
		}
	}
	return categories, nil
}

func (r *categoryRepository) CreateCategory(category *entity.Category) (*entity.Category, error) {
	categoryDB := model.CategoryDB{
		Code: category.Code,
		Name: category.Name,
	}

	if err := r.db.Create(&categoryDB).Error; err != nil {
		return nil, err
	}

	return &entity.Category{
		ID:   categoryDB.ID,
		Code: categoryDB.Code,
		Name: categoryDB.Name,
	}, nil
}
