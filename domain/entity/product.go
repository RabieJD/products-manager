package entity

import "github.com/shopspring/decimal"

// Product represents a product in the domain.
// It includes a unique code, a price, and an optional category.
type Product struct {
	ID       uint
	Code     string
	Price    decimal.Decimal
	Category *Category
	Variants []Variant
}

