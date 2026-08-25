package order

import (
	"context"
	"errors"
	"go-pet-shop/internal/models"
	"go-pet-shop/internal/service"
	"go-pet-shop/internal/storage"
)

type OrderService struct {
	storage OrderStorage
}

type OrderStorage interface {
	CreateOrder(ctx context.Context, order models.Order) (int, error)
	AddOrderItem(ctx context.Context, orderItem models.OrderItem) error
	GetOrderByID(ctx context.Context, id int) (models.Order, error)
	GetOrdersByUserEmail(ctx context.Context, email string) ([]models.Order, error)
	GetOrderItemsByOrderID(ctx context.Context, orderID int) ([]models.OrderItem, error)
}

func NewOrderService(storage OrderStorage) *OrderService {
	return &OrderService{
		storage: storage,
	}
}

func (os *OrderService) CreateOrder(ctx context.Context, order models.Order) (int, error) {

	orderID, err := os.storage.CreateOrder(ctx, order)
	if err != nil {
		return 0, err
	}
	return orderID, nil
}

func (os *OrderService) AddOrderItem(ctx context.Context, orderItem models.OrderItem) error {

	err := os.storage.AddOrderItem(ctx, orderItem)
	if err != nil {
		return err
	}
	return nil
}

func (os *OrderService) GetOrderByID(ctx context.Context, id int) (models.OrderDetail, error) {

	order, err := os.storage.GetOrderByID(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return models.OrderDetail{}, service.ErrNotFound
		}
		return models.OrderDetail{}, err
	}

	items, err := os.storage.GetOrderItemsByOrderID(ctx, id)
	if err != nil {
		return models.OrderDetail{}, err
	}

	return models.OrderDetail{
		Order:      order,
		OrderItems: items,
	}, nil
}

func (os *OrderService) GetOrdersByUserEmail(ctx context.Context, email string) ([]models.Order, error) {

	orders, err := os.storage.GetOrdersByUserEmail(ctx, email)
	if err != nil {
		return []models.Order{}, err
	}
	return orders, nil
}

func (os *OrderService) GetOrderItemsByOrderID(ctx context.Context, orderID int) ([]models.OrderItem, error) {

	orderItems, err := os.storage.GetOrderItemsByOrderID(ctx, orderID)
	if err != nil {
		return []models.OrderItem{}, err
	}
	return orderItems, nil
}
