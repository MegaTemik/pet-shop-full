package service

import (
	"context"
	"errors"
	"go-pet-shop/internal/models"
	"go-pet-shop/internal/storage"
)

type UserService struct {
	storage UserStorage
}

type UserStorage interface {
	CreateUser(ctx context.Context, user models.User) error
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	GetAllUsers(ctx context.Context) ([]models.User, error)
}

func NewUserService(storage UserStorage) *UserService {
	return &UserService{
		storage: storage,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user models.User) error {

	err := s.storage.CreateUser(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (models.User, error) {

	user, err := s.storage.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return models.User{}, ErrNotFound
		}
		return models.User{}, err
	}
	return user, nil
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]models.User, error) {

	users, err := s.storage.GetAllUsers(ctx)
	if err != nil {
		return []models.User{}, err
	}
	return users, nil
}
