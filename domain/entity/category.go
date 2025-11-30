package entity

// Category represents a group of products in the domain.
// It includes a human-readable code and name.
// Categories help organize products into logical groups.
type Category struct {
	ID       uint
	Code     string
	Name     string
	Products []Product
}

