package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/domain/entity"
	domainRepo "github.com/mytheresa/go-hiring-challenge/domain/repository"
	"github.com/mytheresa/go-hiring-challenge/internal/usecase"
	"github.com/stretchr/testify/assert"
)

func TestCategoryHandlerHandleGetCategories(t *testing.T) {
	testCases := []struct {
		name           string
		repo           *mockCategoryRepository
		expectedStatus int
		assertFunc     func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "success",
			repo: &mockCategoryRepository{
				categories: []entity.Category{
					{Code: "CAT001", Name: "Category 1"},
					{Code: "CAT002", Name: "Category 2"},
				},
			},
			expectedStatus: http.StatusOK,
			assertFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp []map[string]string
				err := json.NewDecoder(rec.Body).Decode(&resp)
				assert.NoError(t, err)
				assert.Len(t, resp, 2)
				assert.Equal(t, "CAT001", resp[0]["code"])
			},
		},
		{
			name: "repository_error",
			repo: &mockCategoryRepository{
				getErr: errors.New("db error"),
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			handler := newCategoryHandlerForTest(tc.repo)
			req := httptest.NewRequest(http.MethodGet, "/categories", nil)
			rec := httptest.NewRecorder()

			handler.HandleGetCategories(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			if tc.assertFunc != nil {
				tc.assertFunc(t, rec)
			}
		})
	}
}

func TestCategoryHandlerHandleCreateCategory(t *testing.T) {
	testCases := []struct {
		name           string
		repo           *mockCategoryRepository
		body           string
		expectedStatus int
		assertFunc     func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name:           "success",
			repo:           &mockCategoryRepository{},
			body:           `{"code":"CAT123","name":"Category 123"}`,
			expectedStatus: http.StatusCreated,
			assertFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp map[string]string
				err := json.NewDecoder(rec.Body).Decode(&resp)
				assert.NoError(t, err)
				assert.Equal(t, "CAT123", resp["code"])
				assert.Equal(t, "Category 123", resp["name"])
			},
		},
		{
			name:           "invalid_payload",
			repo:           &mockCategoryRepository{},
			body:           `invalid json`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing_fields",
			repo:           &mockCategoryRepository{},
			body:           `{"code":"","name":" "}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "repository_error",
			repo: &mockCategoryRepository{
				createErr: errors.New("insert failed"),
			},
			body:           `{"code":"CAT001","name":"Category 1"}`,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			handler := newCategoryHandlerForTest(tc.repo)
			req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBufferString(tc.body))
			rec := httptest.NewRecorder()

			handler.HandleCreateCategory(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			if tc.assertFunc != nil && tc.expectedStatus == http.StatusCreated {
				tc.assertFunc(t, rec)
			}
		})
	}
}

func newCategoryHandlerForTest(repo domainRepo.CategoryRepository) *CategoryHandler {
	categoryUsecase := usecase.NewCategoryUsecase(repo)
	return NewCategoryHandler(categoryUsecase)
}

type mockCategoryRepository struct {
	categories []entity.Category
	getErr     error
	createErr  error
}

func (m *mockCategoryRepository) GetAllCategories() ([]entity.Category, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	return m.categories, nil
}

func (m *mockCategoryRepository) CreateCategory(category *entity.Category) (*entity.Category, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}

	return &entity.Category{
		ID:   1,
		Code: category.Code,
		Name: category.Name,
	}, nil
}
