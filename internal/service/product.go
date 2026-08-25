package service

import (
	"context"
	"errors"
	"go-pet-shop/internal/models"
	"go-pet-shop/internal/storage"
)

var (
	ErrNotFound = errors.New("not found")
)

type ProductService struct {
	storage ProductStorage
}

type ProductStorage interface {
	CreateProduct(ctx context.Context, product models.Product) error
	GetProductByID(ctx context.Context, id int) (models.Product, error)
	GetAllProducts(ctx context.Context) ([]models.Product, error)
	UpdateProduct(ctx context.Context, product models.Product) error
	DeleteProduct(ctx context.Context, id int) error
}

func NewProductService(storage ProductStorage) *ProductService {
	return &ProductService{
		storage: storage,
	}
}

func (s *ProductService) CreateProduct(ctx context.Context, product models.Product) error {

	err := s.storage.CreateProduct(ctx, product)
	if err != nil {
		return err
	}
	return nil
}

func (s *ProductService) GetProductByID(ctx context.Context, id int) (models.Product, error) {

	product, err := s.storage.GetProductByID(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return models.Product{}, ErrNotFound
		}
		return models.Product{}, err
	}
	return product, nil
}

func (s *ProductService) GetAllProducts(ctx context.Context) ([]models.Product, error) {

	products, err := s.storage.GetAllProducts(ctx)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, product models.Product) error {

	err := s.storage.UpdateProduct(ctx, product)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id int) error {

	err := s.storage.DeleteProduct(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}
