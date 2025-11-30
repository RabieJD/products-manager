package entity

import "github.com/shopspring/decimal"

// Variant represents a product variant in the domain.
// It includes a unique name, SKU, and an optional price.
// Variants can be used to represent different configurations or options for a product.
type Variant struct {
	ID        uint
	ProductID uint
	Name      string
	SKU       string
	Price     decimal.Decimal
}

