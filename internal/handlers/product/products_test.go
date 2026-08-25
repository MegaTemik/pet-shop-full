package product

import (
	"bytes"
	"errors"
	"fmt"
	"go-pet-shop/internal/handlers/product/mocks"
	"go-pet-shop/internal/models"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/mock"
)

func TestGetAllProducts_Success(t *testing.T) {
	// Мокаем storage — он вернёт один продукт.
	GetAllProductsMock := mocks.NewProductService(t)
	GetAllProductsMock.
		On("GetAllProducts", mock.Anything).
		Return([]models.Product{{ID: 1, Name: "Dog Food"}}, nil).
		Once()

	// Создаем HTTP-запрос GET /products
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	// Создаем хендлер с мок-хранилищем
	handler := New(slog.Default(), GetAllProductsMock)

	// Вызываем метод GetAllProducts
	handler.GetAllProducts(w, req)

	// Проверяем HTTP-код
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}
func TestGetAllProducts_Fail(t *testing.T) {
	// Мокаем storage — он будет возвращать ошибку
	GetAllProductsMock := mocks.NewProductService(t)
	GetAllProductsMock.
		On("GetAllProducts", mock.Anything).
		Return(nil, errors.New("DB error")).
		Once()

	// Создаем HTTP-запрос GET /products
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	// Создаем хендлер с мок-хранилищем
	handler := New(slog.Default(), GetAllProductsMock)

	// Вызываем метод GetAllProducts
	handler.GetAllProducts(w, req)

	// Проверяем HTTP-код
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// =======================
// Create Product
// =======================

// Исправить в v1, добавить once
func TestCreateProduct_Success(t *testing.T) {
	// Мокаем storage - он будет возвращать ID
	CreateProductMock := mocks.NewProductService(t)
	CreateProductMock.
		On("CreateProduct", mock.Anything, mock.Anything).
		Return(nil).
		Once()

	// Создаем HTTP-запрос POST /products
	input := "{\"name\": \"Dog Food\", \"price\": 19.99, \"stock\": 10}"
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader([]byte(input)))
	w := httptest.NewRecorder()

	// Создаем хендлер с мок-хранилищем
	handler := New(slog.Default(), CreateProductMock)

	// Вызываем метод CreateProduct
	handler.CreateProduct(w, req)

	// Проверяем HTTP-код
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestCreateProduct_BadRequest(t *testing.T) {
	// Мокаем storage - он ничего не будет возвращать
	CreateProductMock := mocks.NewProductService(t)

	// Создаем HTTP запрос POST /products
	input := "{\"name\": \"Cat Food\", \"price\": 50.00, \"stock\": 100"
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader([]byte(input)))
	w := httptest.NewRecorder()

	// Создаем хендлер с мок-хранилищем
	handler := New(slog.Default(), CreateProductMock)

	// Вызываем метод CreateProduct
	handler.CreateProduct(w, req)

	// Проверяем HTTP-код
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestCreateProduct_Fail(t *testing.T) {
	// Мокаем storage - он будет возвращать ошибку
	CreateProductMock := mocks.NewProductService(t)
	CreateProductMock.
		On("CreateProduct", mock.Anything, mock.Anything).
		Return(errors.New("server error")).
		Once()

	// Создаем HTTP запрос POST /products
	input := "{\"name\": \"Cat Food\", \"price\": 50.00, \"stock\": 100}"
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader([]byte(input)))
	w := httptest.NewRecorder()

	// Создаем хендлер с мок-хранилищем
	handler := New(slog.Default(), CreateProductMock)

	// Вызываем метод CreateProduct
	handler.CreateProduct(w, req)

	// Проверяем HTTP-код
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

// =======================
// Update Product
// =======================

func TestUpdateProduct_Success(t *testing.T) {
	// Мокаем storage - он будет возвращать nil
	UpdateProductMock := mocks.NewProductService(t)
	UpdateProductMock.
		On("UpdateProduct", mock.Anything, mock.Anything).
		Return(nil).
		Once()

	// Создаем роутер
	r := chi.NewRouter()

	// Создаем HTTP-запрос PUT /products/{id}
	id := 1
	input := "{\"name\": \"Fox food\", \"price\": 150, \"stock\": 200}"
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/products/%d", id), bytes.NewReader([]byte(input)))
	w := httptest.NewRecorder()

	// Создаем хендлер с мок-хранилищем
	handler := New(slog.Default(), UpdateProductMock)

	// Вызываем метод UpdateProduct
	r.Put("/products/{id}", handler.UpdateProduct)
	r.ServeHTTP(w, req)

	// Проверяем HTTP-код
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestUpdateProduct_BadRequest(t *testing.T) {
	// Мокаем storage - он ничего не будет возвращать
	UpdateProductMock := mocks.NewProductService(t)

	// Создаем роутер
	r := chi.NewRouter()

	// Создаем HTTP запрос PUT /products/{id}
	id := 1
	input := "{\"name\": \"Fox food\", \"price\": 150, \"stock\": 200"
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/products/%d", id), bytes.NewReader([]byte(input)))
	w := httptest.NewRecorder()

	// Создаем хендлер с мок-хранилищем
	handler := New(slog.Default(), UpdateProductMock)

	// Вызываем метод UpdateProduct
	r.Put("/products/{id}", handler.UpdateProduct)
	r.ServeHTTP(w, req)

	// Проверяем HTTP код
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestUpdateProduct_Fail(t *testing.T) {
	//Мокаем storage - он будет возвращать ошибку
	UpdateProductMock := mocks.NewProductService(t)
	UpdateProductMock.
		On("UpdateProduct", mock.Anything, mock.Anything).
		Return(errors.New("server error")).
		Once()

	// Создаем роутер
	r := chi.NewRouter()

	// Создаем HTTP Запрос PUT /products/{id}
	id := 1
	input := "{\"name\": \"Fox food\", \"price\": 150, \"stock\": 200}"
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/products/%d", id), bytes.NewReader([]byte(input)))
	w := httptest.NewRecorder()

	// Создаем хендлер с мок-хранилищем
	handler := New(slog.Default(), UpdateProductMock)

	// Вызываем метод UpdateProduct
	r.Put("/products/{id}", handler.UpdateProduct)
	r.ServeHTTP(w, req)

	// Проверяем HTTP код
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

// =======================
// Delete Product
// =======================

func TestDeleteProduct_Success(t *testing.T) {
	// Мокаем storage - он будет возвращать nil
	DeleteProductMock := mocks.NewProductService(t)
	DeleteProductMock.
		On("DeleteProduct", mock.Anything, mock.Anything).
		Return(nil).
		Once()

	// Создаем роутер
	r := chi.NewRouter()

	// Создаем HTTP Запрос DELETE /products/{id}
	id := 1
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/products/%d", id), nil)
	w := httptest.NewRecorder()

	// Создаем хендлер с мок-хранилищем
	handler := New(slog.Default(), DeleteProductMock)

	// Вызываем метод DeleteProduct
	r.Delete("/products/{id}", handler.DeleteProduct)
	r.ServeHTTP(w, req)

	// Проверяем HTTP код
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestDeleteProduct_BadRequest(t *testing.T) {
	// Мокаем storage - он ничего не будет возвращать
	DeleteProductMock := mocks.NewProductService(t)

	// Создаем HTTP запрос DELETE /products/{id}
	req := httptest.NewRequest(http.MethodDelete, "/products/", nil)
	w := httptest.NewRecorder()

	// Создаем хендлер с мок-хранилищем
	handler := New(slog.Default(), DeleteProductMock)

	// Вызываем метод DeleteProduct
	handler.DeleteProduct(w, req)

	// Проверяем HTTP-код
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestDeleteProduct_Fail(t *testing.T) {
	// Мокаем storage - он будет возвращать ошибку
	DeleteProductMock := mocks.NewProductService(t)
	DeleteProductMock.
		On("DeleteProduct", mock.Anything, mock.Anything).
		Return(errors.New("server error")).
		Once()

	// Создаем роутер
	r := chi.NewRouter()

	// Создаем HTTP Запрос DELETE /products/{id}
	id := 1
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/products/%d", id), nil)
	w := httptest.NewRecorder()

	// Создаем хендлер с мок-хранилищем
	handler := New(slog.Default(), DeleteProductMock)

	// Вызываем метод DeleteProduct
	r.Delete("/products/{id}", handler.DeleteProduct)
	r.ServeHTTP(w, req)

	// Проверяем HTTP код
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

func TestGetProductByID_Success(t *testing.T) {
	// Мокаем storage - он будет возвращать продукт
	GetProductByIDMock := mocks.NewProductService(t)
	GetProductByIDMock.
		On("GetProductByID", mock.Anything, mock.Anything).
		Return(models.Product{ID: 1, Name: "Dog Food"}, nil).
		Once()

	// Создаем роутер
	r := chi.NewRouter()

	// Создаем HTTP Запрос GET /products/{id}
	id := 1
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/products/%d", id), nil)
	w := httptest.NewRecorder()

	// Создаем хендлер с мок-хранилищем
	handler := New(slog.Default(), GetProductByIDMock)

	// Вызываем метод GetProductByID
	r.Get("/products/{id}", handler.GetProductByID)
	r.ServeHTTP(w, req)

	// Проверяем HTTP код
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestGetProductByID_Fail(t *testing.T) {
	// Мокаем storage - он будет возвращать ошибку
	GetProductByIDMock := mocks.NewProductService(t)
	GetProductByIDMock.
		On("GetProductByID", mock.Anything, mock.Anything).
		Return(models.Product{}, errors.New("DB error")).
		Once()

	// Создаем роутер
	r := chi.NewRouter()

	// Создаем HTTP Запрос GET /products/{id}
	id := 1
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/products/%d", id), nil)
	w := httptest.NewRecorder()

	// Создаем хендлер с мок-хранилищем
	handler := New(slog.Default(), GetProductByIDMock)

	// Вызываем метод GetProductByID
	r.Get("/products/{id}", handler.GetProductByID)
	r.ServeHTTP(w, req)

	// Проверяем HTTP код
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

// Добавить к v1 id := nooo
func TestGetProductByID_BadRequest(t *testing.T) {
	// Мокаем storage - он ничего не будет возвращать
	GetProductByIDMock := mocks.NewProductService(t)

	// Создаем HTTP Запрос GET /products/{id}
	id := "not-a-number"
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/products/%s", id), nil)
	w := httptest.NewRecorder()

	// Создаем хендлер с мок-хранилищем
	handler := New(slog.Default(), GetProductByIDMock)

	// Вызываем метод GetProductByID
	handler.GetProductByID(w, req)

	// Проверяем HTTP код
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestGetPopularProducts_Success(t *testing.T) {
	GetPopularProductsMock := mocks.NewProductService(t)
	GetPopularProductsMock.
		On("GetPopularProducts", mock.Anything).
		Return([]models.PopularProduct{
			{
				Product: models.Product{
					ID:    1,
					Name:  "Яблоко",
					Price: 100,
					Stock: 1000,
				},
				Count: 500,
			},
		}, nil).
		Once()

	req := httptest.NewRequest(http.MethodGet, "/products/popular", nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), GetPopularProductsMock)
	handler.GetPopularProducts(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestGetPopularProducts_Fail(t *testing.T) {
	GetPopularProductsMock := mocks.NewProductService(t)
	GetPopularProductsMock.
		On("GetPopularProducts", mock.Anything).
		Return([]models.PopularProduct{
			{
				Product: models.Product{
					ID:    1,
					Name:  "Яблоко",
					Price: 100,
					Stock: 1000,
				},
				Count: 500,
			},
		}, errors.New("server error")).
		Once()

	req := httptest.NewRequest(http.MethodGet, "/products/popular", nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), GetPopularProductsMock)
	handler.GetPopularProducts(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}
