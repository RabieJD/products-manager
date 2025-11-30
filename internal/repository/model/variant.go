package model

import (
	"github.com/shopspring/decimal"
)

// VariantDB represents the database model for product variants.
type VariantDB struct {
	ID        uint            `gorm:"primaryKey"`
	ProductID uint            `gorm:"not null"`
	Name      string          `gorm:"not null"`
	SKU       string          `gorm:"uniqueIndex;not null"`
	Price     decimal.Decimal `gorm:"type:decimal(10,2);null"`
}

func (VariantDB) TableName() string {
	return "product_variants"
}


