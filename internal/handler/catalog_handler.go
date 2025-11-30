package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/domain/repository"
	"github.com/mytheresa/go-hiring-challenge/internal/dto"
	"github.com/mytheresa/go-hiring-challenge/internal/usecase"
	"github.com/shopspring/decimal"
)

// CatalogHandler handles HTTP requests for the catalog.
type CatalogHandler struct {
	productUsecase *usecase.ProductUsecase
}

// NewCatalogHandler creates a new catalog handler instance.
func NewCatalogHandler(productUsecase *usecase.ProductUsecase) *CatalogHandler {
	return &CatalogHandler{
		productUsecase: productUsecase,
	}
}

// HandleGet handles GET /catalog requests.
func (h *CatalogHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	offset := h.parseIntQueryParam(r, "offset", 0)
	limit := h.parseIntQueryParam(r, "limit", 10)

	// Validate and constrain limit
	if limit < 1 {
		limit = 1
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	// Parse filter parameters
	filter := &repository.ProductFilter{}

	if categoryCode := r.URL.Query().Get("category"); categoryCode != "" {
		filter.CategoryCode = &categoryCode
	}

	if priceLessThanStr := r.URL.Query().Get("price_less_than"); priceLessThanStr != "" {
		priceLessThan, err := decimal.NewFromString(priceLessThanStr)
		if err == nil {
			filter.PriceLessThan = &priceLessThan
		}
	}

	// Build pagination params
	pagination := &repository.PaginationParams{
		Offset: offset,
		Limit:  limit,
	}

	// Get products from usecase
	result, err := h.productUsecase.GetProducts(filter, pagination)
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Map domain entities to DTOs
	productDTOs := make([]dto.ProductResponse, len(result.Products))
	for i, p := range result.Products {
		var categoryDTO *dto.CategoryResponse
		if p.Category != nil {
			categoryDTO = &dto.CategoryResponse{
				Code: p.Category.Code,
				Name: p.Category.Name,
			}
		}

		productDTOs[i] = dto.ProductResponse{
			Code:     p.Code,
			Price:    p.Price.InexactFloat64(),
			Category: categoryDTO,
		}
	}

	api.OKResponse(w, dto.CatalogResponse{
		Products: productDTOs,
		Total:    result.Total,
		Offset:   offset,
		Limit:    limit,
	})
}

// HandleGetProduct handles GET /catalog/{code} requests.
func (h *CatalogHandler) HandleGetProduct(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimSpace(r.PathValue("code"))
	if code == "" {
		api.ErrorResponse(w, http.StatusBadRequest, "product code is required")
		return
	}

	product, err := h.productUsecase.GetProductByCode(code)
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if product == nil {
		api.ErrorResponse(w, http.StatusNotFound, "product not found")
		return
	}

	var categoryDTO *dto.CategoryResponse
	if product.Category != nil {
		categoryDTO = &dto.CategoryResponse{
			Code: product.Category.Code,
			Name: product.Category.Name,
		}
	}

	variantDTOs := make([]dto.VariantResponse, len(product.Variants))
	for i, v := range product.Variants {
		price := v.Price
		if price.Equal(decimal.Zero) {
			price = product.Price
		}

		variantDTOs[i] = dto.VariantResponse{
			Name:  v.Name,
			SKU:   v.SKU,
			Price: price.InexactFloat64(),
		}
	}

	api.OKResponse(w, dto.ProductDetailsResponse{
		Code:     product.Code,
		Price:    product.Price.InexactFloat64(),
		Category: categoryDTO,
		Variants: variantDTOs,
	})
}

// parseIntQueryParam parses an integer query parameter with a default value.
func (h *CatalogHandler) parseIntQueryParam(r *http.Request, key string, defaultValue int) int {
	valueStr := r.URL.Query().Get(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}
