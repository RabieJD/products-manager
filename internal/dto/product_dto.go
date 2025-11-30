package dto

// CategoryResponse represents a category in the API response.
type CategoryResponse struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// CategoryRequest represents incoming payload to create a category.
type CategoryRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// ProductResponse represents a product in the API response.
type ProductResponse struct {
	Code     string            `json:"code"`
	Price    float64           `json:"price"`
	Category *CategoryResponse `json:"category,omitempty"`
}

// VariantResponse represents a product variant in the API response.
type VariantResponse struct {
	Name  string  `json:"name"`
	SKU   string  `json:"sku"`
	Price float64 `json:"price"`
}

// ProductDetailsResponse represents the detailed response for a single product.
type ProductDetailsResponse struct {
	Code     string            `json:"code"`
	Price    float64           `json:"price"`
	Category *CategoryResponse `json:"category,omitempty"`
	Variants []VariantResponse `json:"variants"`
}

// CatalogResponse represents the catalog API response.
type CatalogResponse struct {
	Products []ProductResponse `json:"products"`
	Total    int64             `json:"total"`
	Offset   int               `json:"offset"`
	Limit    int               `json:"limit"`
}
