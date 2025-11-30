package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/domain/entity"
	"github.com/mytheresa/go-hiring-challenge/domain/repository"
	"github.com/mytheresa/go-hiring-challenge/internal/dto"
	"github.com/mytheresa/go-hiring-challenge/internal/usecase"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCatalogHandlerHandleGetProduct(t *testing.T) {
	productWithVariants := &entity.Product{
		Code:  "PROD001",
		Price: decimal.NewFromFloat(10.99),
		Category: &entity.Category{
			Code: "CAT001",
			Name: "Category 1",
		},
		Variants: []entity.Variant{
			{Name: "Variant A", SKU: "SKU001A", Price: decimal.Zero},
			{Name: "Variant B", SKU: "SKU001B", Price: decimal.NewFromFloat(12.49)},
		},
	}

	testCases := []struct {
		name           string
		product        *entity.Product
		productErr     error
		productCode    string
		expectedStatus int
		assertFunc     func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name:        "success_with_variants_and_category",
			product:     productWithVariants,
			productCode: "PROD001",
			assertFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response dto.ProductDetailsResponse
				err := json.NewDecoder(rec.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, productWithVariants.Code, response.Code)
				assert.Equal(t, productWithVariants.Price.InexactFloat64(), response.Price)
				assert.NotNil(t, response.Category)
				assert.Equal(t, 2, len(response.Variants))
				assert.Equal(t, productWithVariants.Price.InexactFloat64(), response.Variants[0].Price)
				assert.Equal(t, productWithVariants.Variants[1].Price.InexactFloat64(), response.Variants[1].Price)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "not_found",
			product:        nil,
			productCode:    "UNKNOWN",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "repository_error",
			productErr:     errors.New("database error"),
			productCode:    "PROD001",
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			handler := newCatalogHandlerForTest(tc.product, tc.productErr)
			req := httptest.NewRequest(http.MethodGet, "/catalog/"+tc.productCode, nil)
			req.SetPathValue("code", tc.productCode)
			rec := httptest.NewRecorder()

			handler.HandleGetProduct(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			if tc.assertFunc != nil && tc.expectedStatus == http.StatusOK {
				tc.assertFunc(t, rec)
			}
		})
	}
}

func TestCatalogHandlerHandleGet(t *testing.T) {
	testCases := []struct {
		name           string
		repo           *mockProductRepository
		url            string
		expectedStatus int
		assertFunc     func(t *testing.T, rec *httptest.ResponseRecorder, repo *mockProductRepository)
	}{
		{
			name: "success_with_pagination",
			repo: &mockProductRepository{
				listResult: &repository.PaginationResult{
					Products: []entity.Product{
						{
							Code:  "PROD001",
							Price: decimal.NewFromFloat(10.5),
							Category: &entity.Category{
								Code: "CAT001",
								Name: "Category 1",
							},
						},
					},
					Total: 1,
				},
			},
			url:            "/catalog?limit=5&offset=1",
			expectedStatus: http.StatusOK,
			assertFunc: func(t *testing.T, rec *httptest.ResponseRecorder, repo *mockProductRepository) {
				var response dto.CatalogResponse
				err := json.NewDecoder(rec.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, int64(1), response.Total)
				assert.Equal(t, 1, response.Offset)
				assert.Equal(t, 5, response.Limit)
				assert.Len(t, response.Products, 1)
				assert.Equal(t, "PROD001", response.Products[0].Code)
				assert.NotNil(t, response.Products[0].Category)
			},
		},
		{
			name: "repository_error",
			repo: &mockProductRepository{
				listErr: errors.New("query failed"),
			},
			url:            "/catalog",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "applies_filters_and_clamps_pagination",
			repo: &mockProductRepository{
				listResult: &repository.PaginationResult{},
			},
			url:            "/catalog?category=CAT123&price_less_than=9.99&offset=-5&limit=150",
			expectedStatus: http.StatusOK,
			assertFunc: func(t *testing.T, rec *httptest.ResponseRecorder, repo *mockProductRepository) {
				assert.Equal(t, http.StatusOK, rec.Code)
				if assert.NotNil(t, repo.lastFilter) {
					if assert.NotNil(t, repo.lastFilter.CategoryCode) {
						assert.Equal(t, "CAT123", *repo.lastFilter.CategoryCode)
					}
					if assert.NotNil(t, repo.lastFilter.PriceLessThan) {
						assert.Equal(t, decimal.RequireFromString("9.99"), *repo.lastFilter.PriceLessThan)
					}
				}
				if assert.NotNil(t, repo.lastPagination) {
					assert.Equal(t, 0, repo.lastPagination.Offset)
					assert.Equal(t, 100, repo.lastPagination.Limit)
				}
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			handler := newCatalogHandlerWithRepo(tc.repo)

			req := httptest.NewRequest(http.MethodGet, tc.url, nil)
			rec := httptest.NewRecorder()

			handler.HandleGet(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			if tc.assertFunc != nil {
				tc.assertFunc(t, rec, tc.repo)
			}
		})
	}
}

// newCatalogHandlerForTest creates a catalog handler with a mocked repository for product lookups.
func newCatalogHandlerForTest(product *entity.Product, productErr error) *CatalogHandler {
	repo := &mockProductRepository{
		product:    product,
		productErr: productErr,
	}

	return newCatalogHandlerWithRepo(repo)
}

func newCatalogHandlerWithRepo(repo repository.ProductRepository) *CatalogHandler {
	productUsecase := usecase.NewProductUsecase(repo)
	return NewCatalogHandler(productUsecase)
}

type mockProductRepository struct {
	product    *entity.Product
	productErr error

	listResult *repository.PaginationResult
	listErr    error

	lastFilter     *repository.ProductFilter
	lastPagination *repository.PaginationParams
}

func (m *mockProductRepository) GetProducts(filter *repository.ProductFilter, pagination *repository.PaginationParams) (*repository.PaginationResult, error) {
	if filter != nil {
		m.lastFilter = &repository.ProductFilter{}
		if filter.CategoryCode != nil {
			code := *filter.CategoryCode
			m.lastFilter.CategoryCode = &code
		}
		if filter.PriceLessThan != nil {
			price := *filter.PriceLessThan
			m.lastFilter.PriceLessThan = &price
		}
	} else {
		m.lastFilter = nil
	}

	if pagination != nil {
		cp := *pagination
		m.lastPagination = &cp
	} else {
		m.lastPagination = nil
	}

	if m.listErr != nil {
		return nil, m.listErr
	}
	if m.listResult != nil {
		return m.listResult, nil
	}
	return &repository.PaginationResult{}, nil
}

func (m *mockProductRepository) GetProductByCode(code string) (*entity.Product, error) {
	return m.product, m.productErr
}

func (m *mockProductRepository) GetProductByID(id uint) (*entity.Product, error) {
	return nil, nil
}
