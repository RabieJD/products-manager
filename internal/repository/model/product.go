package model

import (
	"github.com/shopspring/decimal"
)

// ProductDB represents the database model for products.
// This is separate from the domain entity to maintain clean architecture.
type ProductDB struct {
	ID         uint            `gorm:"primaryKey"`
	Code       string          `gorm:"uniqueIndex;not null"`
	Price      decimal.Decimal `gorm:"type:decimal(10,2);not null"`
	CategoryID *uint           `gorm:"index"`
	Category   *CategoryDB     `gorm:"foreignKey:CategoryID"`
	Variants   []VariantDB      `gorm:"foreignKey:ProductID"`
}

func (ProductDB) TableName() string {
	return "products"
}


