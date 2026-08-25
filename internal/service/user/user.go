package product

import (
	"context"
	"go-pet-shop/internal/models"
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

func (us *UserService) CreateUser(ctx context.Context, user models.User) error {

	err := us.storage.CreateUser(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (us *UserService) GetUserByEmail(ctx context.Context, email string) (models.User, error) {

	user, err := us.storage.GetUserByEmail(ctx, email)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (us *UserService) GetAllUsers(ctx context.Context) ([]models.User, error) {

	users, err := us.storage.GetAllUsers(ctx)
	if err != nil {
		return []models.User{}, err
	}
	return users, nil
}
