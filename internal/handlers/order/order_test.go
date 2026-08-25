package order

import (
	"bytes"
	"errors"
	"fmt"
	"go-pet-shop/internal/handlers/order/mocks"
	"go-pet-shop/internal/models"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/mock"
)

func TestCreateOrder_Success(t *testing.T) {
	CreateOrderMock := mocks.NewOrderService(t)
	CreateOrderMock.
		On("CreateOrder", mock.Anything, mock.Anything).
		Return(1, nil).
		Once()

	input := "{\"user_email\": \"john.doe@example.com\", \"total_price\": 100.0}"
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader([]byte(input)))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), CreateOrderMock)

	handler.CreateOrder(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestCreateOrder_Fail(t *testing.T) {
	CreateOrderMock := mocks.NewOrderService(t)
	CreateOrderMock.
		On("CreateOrder", mock.Anything, mock.Anything).
		Return(0, errors.New("server error")).
		Once()

	input := "{\"user_email\": \"john.doe@example.com\", \"total_price\": 100.0}"
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader([]byte(input)))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), CreateOrderMock)

	handler.CreateOrder(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

func TestCreateOrder_BadRequest(t *testing.T) {
	CreateOrderMock := mocks.NewOrderService(t)

	input := "{\"user_email\": \"\", \"total_price\": 100.0}" // Empty email
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader([]byte(input)))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), CreateOrderMock)

	handler.CreateOrder(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestAddOrderItem_Success(t *testing.T) {
	AddOrderItemMock := mocks.NewOrderService(t)
	AddOrderItemMock.
		On("AddOrderItem", mock.Anything, mock.Anything).
		Return(nil).
		Once()

	r := chi.NewRouter()

	id := 1
	input := "{\"order_id\": 1, \"product_id\": 2, \"quantity\": 3}"
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/orders/%d/items", id), bytes.NewReader([]byte(input)))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), AddOrderItemMock)

	r.Post("/orders/{id}/items", handler.AddOrderItem)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestAddOrderItem_Fail(t *testing.T) {
	AddOrderItemMock := mocks.NewOrderService(t)
	AddOrderItemMock.
		On("AddOrderItem", mock.Anything, mock.Anything).
		Return(errors.New("server error")).
		Once()

	r := chi.NewRouter()

	id := 1
	input := "{\"order_id\": 1, \"product_id\": 2, \"quantity\": 3}"
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/orders/%d/items", id), bytes.NewReader([]byte(input)))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), AddOrderItemMock)

	r.Post("/orders/{id}/items", handler.AddOrderItem)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

func TestAddOrderItem_BadRequest(t *testing.T) {
	AddOrderItemMock := mocks.NewOrderService(t)

	r := chi.NewRouter()

	id := 1
	input := "{\"order_id\": 1, \"product_id\": 2"
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/orders/%d/items", id), bytes.NewReader([]byte(input)))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), AddOrderItemMock)

	r.Post("/orders/{id}/items", handler.AddOrderItem)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestGetOrderByID_Success(t *testing.T) {
	GetOrderByIDMock := mocks.NewOrderService(t)
	GetOrderByIDMock.
		On("GetOrderByID", mock.Anything, mock.Anything).
		Return(models.Order{
			ID:         1,
			UserEmail:  "Bob",
			TotalPrice: 111.1,
			CreatedAt:  time.Now(),
		}, nil)

	r := chi.NewRouter()
	id := 1
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/orders/%d", id), nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), GetOrderByIDMock)

	r.Get("/orders/{id}", handler.GetOrderByID)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestGetOrderByID_Fail(t *testing.T) {
	GetOrderByIDMock := mocks.NewOrderService(t)
	GetOrderByIDMock.
		On("GetOrderByID", mock.Anything, mock.Anything).
		Return(models.Order{}, errors.New("server error")).
		Once()

	r := chi.NewRouter()

	id := 1
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/orders/%d", id), nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), GetOrderByIDMock)

	r.Get("/orders/{id}", handler.GetOrderByID)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

func TestGetOrderById_BadRequest(t *testing.T) {
	GetOrderByIDMock := mocks.NewOrderService(t)

	r := chi.NewRouter()

	id := "not-a-number"
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/orders/%s", id), nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), GetOrderByIDMock)

	r.Get("/orders/{id}", handler.GetOrderByID)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestGetOrdersByUserEmail_Success(t *testing.T) {
	GetOrdersByUserEmailMock := mocks.NewOrderService(t)
	GetOrdersByUserEmailMock.
		On("GetOrdersByUserEmail", mock.Anything, mock.Anything).
		Return([]models.Order{
			{
				ID:         1,
				UserEmail:  "Bob@gmail.com",
				TotalPrice: 111.1,
				CreatedAt:  time.Now(),
			},
		}, nil)

	email := "Bob@gmail.com"
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/users/orders?email=%s", email), nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), GetOrdersByUserEmailMock)

	handler.GetOrdersByUserEmail(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestGetOrdersByUserEmail_Fail(t *testing.T) {
	GetOrdersByUserEmailMock := mocks.NewOrderService(t)
	GetOrdersByUserEmailMock.
		On("GetOrdersByUserEmail", mock.Anything, mock.Anything).
		Return([]models.Order{}, errors.New("server error")).
		Once()

	email := "Bob@gmail.com"
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/users/orders?email=%s", email), nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), GetOrdersByUserEmailMock)

	handler.GetOrdersByUserEmail(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

func TestGetOrdersByUserEmail_BadRequest(t *testing.T) {
	GetOrdersByUserEmailMock := mocks.NewOrderService(t)

	email := ""
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/users/orders?email=%s", email), nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), GetOrdersByUserEmailMock)

	handler.GetOrdersByUserEmail(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}
