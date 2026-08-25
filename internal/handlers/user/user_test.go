package user

import (
	"bytes"
	"errors"
	"fmt"
	"go-pet-shop/internal/handlers/user/mocks"
	"go-pet-shop/internal/models"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/mock"
)

func TestGetAllUsers_Success(t *testing.T) {
	// Мокаем storage - он будет возвращать список пользователей
	GetAllUsersMock := mocks.NewUserService(t)
	GetAllUsersMock.
		On("GetAllUsers", mock.Anything).
		Return([]models.User{
			{ID: 1, Name: "John Doe", Email: "john.doe@example.com"},
			{ID: 2, Name: "Jane Smith", Email: "jane.smith@example.com"},
		}, nil).
		Once()

	// Создаем HTTP-запрос GET /users
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	// Создаем хендлер с мок-хранилищем
	handler := New(slog.Default(), GetAllUsersMock)

	// Вызываем метод GetAllUsers
	handler.GetAllUsers(w, req)

	// Проверяем HTTP-код
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestGetAllUsers_Fail(t *testing.T) {
	// Мокаем storage - он будет возвращать ошибку
	GetAllUsersMock := mocks.NewUserService(t)
	GetAllUsersMock.
		On("GetAllUsers", mock.Anything).
		Return(nil, errors.New("server error")).
		Once()

	// Создаем HTTP-запрос GET /users
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	// Создаем хендлер с мок-хранилищем
	handler := New(slog.Default(), GetAllUsersMock)

	// Вызываем метод GetAllUsers
	handler.GetAllUsers(w, req)

	// Проверяем HTTP-код
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

func TestGetUserByEmail_Success(t *testing.T) {
	// Мокаем storage - он будет возвращать пользователя
	GetUserByEmailMock := mocks.NewUserService(t)
	GetUserByEmailMock.
		On("GetUserByEmail", mock.Anything, "john.doe@example.com").
		Return(models.User{ID: 1, Name: "John Doe", Email: "john.doe@example.com"}, nil).
		Once()

	// Создаем роутер
	r := chi.NewRouter()

	// Создаем HTTP запрос GET /products/{email}
	email := "john.doe@example.com"
	req := httptest.NewRequest(http.MethodGet, "/users/"+email, nil)
	w := httptest.NewRecorder()

	// Создаем хендлер с мок-хранилищем
	handler := New(slog.Default(), GetUserByEmailMock)

	// Вызываем метод GetUserByEmail
	r.Get("/users/{email}", handler.GetUserByEmail)
	r.ServeHTTP(w, req)

	// Проверяем HTTP-код
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestGetUserByEmail_Fail(t *testing.T) {
	// Мокаем storage - он будет возвращать ошибку
	GetUserByEmailMock := mocks.NewUserService(t)
	GetUserByEmailMock.
		On("GetUserByEmail", mock.Anything, mock.Anything).
		Return(models.User{}, errors.New("DB error")).
		Once()

	// Создаем роутер
	r := chi.NewRouter()

	// Создаем HTTP запрос GET /products/{email}
	email := "john.doe@example.com"
	req := httptest.NewRequest(http.MethodGet, "/users/"+email, nil)
	w := httptest.NewRecorder()

	// Создаем хендлер с мок-хранилищем
	handler := New(slog.Default(), GetUserByEmailMock)

	// Вызываем метод GetUserByEmail
	r.Get("/users/{email}", handler.GetUserByEmail)
	r.ServeHTTP(w, req)

	// Проверяем HTTP-код
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

func TestGetUserByEmail_BadRequest(t *testing.T) {
	// Мокаем storage - он ничего не будет возвращать
	GetUserByEmailMock := mocks.NewUserService(t)

	email := ""
	// Создаем HTTP запрос GET /products/{email}
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/users/%s", email), nil)
	w := httptest.NewRecorder()

	// Создаем хендлер с мок-хранилищем
	handler := New(slog.Default(), GetUserByEmailMock)

	// Вызываем метод GetUserByEmail
	handler.GetUserByEmail(w, req)

	// Проверяем HTTP-код
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestCreateUser_Success(t *testing.T) {
	// Мокаем storage - он ничего не будет возвращать
	CreateUserMock := mocks.NewUserService(t)
	CreateUserMock.
		On("CreateUser", mock.Anything, mock.Anything).
		Return(nil).
		Once()

	// Создаем HTTP запрос POST /users
	input := "{\"name\": \"John Doe\", \"email\": \"john.doe@example.com\"}"
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte(input)))
	w := httptest.NewRecorder()

	// Создаем хендлер с мок-хранилищем
	handler := New(slog.Default(), CreateUserMock)

	// Вызываем метод CreateUser
	handler.CreateUser(w, req)

	// Проверяем HTTP-код
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestCreateUser_Fail(t *testing.T) {
	// Мокаем storage - он будет возвращать ошибку
	CreateUserMock := mocks.NewUserService(t)
	CreateUserMock.
		On("CreateUser", mock.Anything, mock.Anything).
		Return(errors.New("server error")).
		Once()

	// Создаем HTTP запрос POST /users
	input := "{\"name\": \"John Doe\", \"email\": \"john.doe@example.com\"}"
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte(input)))
	w := httptest.NewRecorder()

	// Создаем хендлер с мок-хранилищем
	handler := New(slog.Default(), CreateUserMock)

	// Вызываем метод CreateUser
	handler.CreateUser(w, req)

	// Проверяем HTTP-код
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

func TestCreateUser_BadRequest(t *testing.T) {
	// Мокаем storage - он ничего не будет возвращать
	CreateUserMock := mocks.NewUserService(t)

	// Создаем HTTP запрос POST /users
	input := "{\"name\": \"\", \"email\": \"\"}"
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte(input)))
	w := httptest.NewRecorder()

	// Создаем хендлер с мок-хранилищем
	handler := New(slog.Default(), CreateUserMock)

	// Вызываем метод CreateUser
	handler.CreateUser(w, req)

	// Проверяем HTTP-код
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}
