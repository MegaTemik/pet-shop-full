package service

import (
	"context"
	"errors"
	"go-pet-shop/internal/models"
	"go-pet-shop/internal/service"
	"go-pet-shop/internal/storage"
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

func (ps *ProductService) CreateProduct(ctx context.Context, product models.Product) error {

	err := ps.storage.CreateProduct(ctx, product)
	if err != nil {
		return err
	}
	return nil
}

func (ps *ProductService) GetProductByID(ctx context.Context, id int) (models.Product, error) {

	product, err := ps.storage.GetProductByID(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return models.Product{}, service.ErrNotFound
		}
		return models.Product{}, err
	}
	return product, nil
}

func (ps *ProductService) GetAllProducts(ctx context.Context) ([]models.Product, error) {

	products, err := ps.storage.GetAllProducts(ctx)
	if err != nil {
		return []models.Product{}, err
	}
	return products, nil
}

func (ps *ProductService) UpdateProduct(ctx context.Context, product models.Product) error {

	err := ps.storage.UpdateProduct(ctx, product)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return service.ErrNotFound
		}
		return err
	}
	return nil
}

func (ps *ProductService) DeleteProduct(ctx context.Context, id int) error {

	err := ps.storage.DeleteProduct(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return service.ErrNotFound
		}
		return err
	}
	return nil
}
