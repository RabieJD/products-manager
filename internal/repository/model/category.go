package model

// CategoryDB represents the database model for categories.
type CategoryDB struct {
	ID       uint        `gorm:"primaryKey"`
	Code     string      `gorm:"uniqueIndex;not null"`
	Name     string      `gorm:"not null"`
	Products []ProductDB `gorm:"foreignKey:CategoryID"`
}

func (CategoryDB) TableName() string {
	return "product_categories"
}

