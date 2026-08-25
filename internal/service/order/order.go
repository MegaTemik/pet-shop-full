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
	PlaceOrder(ctx context.Context, userEmail string, items []models.OrderItem) (orderID int, err error)
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

func (os *OrderService) GetOrderByID(ctx context.Context, id int) (models.Order, error) {

	order, err := os.storage.GetOrderByID(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return models.Order{}, service.ErrNotFound
		}
		return models.Order{}, err
	}

	items, err := os.storage.GetOrderItemsByOrderID(ctx, id)
	if err != nil {
		return models.Order{}, err
	}
	order.Items = items

	return order, nil
}

func (os *OrderService) GetOrdersByUserEmail(ctx context.Context, email string) ([]models.Order, error) {

	orders, err := os.storage.GetOrdersByUserEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (os *OrderService) GetOrderItemsByOrderID(ctx context.Context, orderID int) ([]models.OrderItem, error) {

	items, err := os.storage.GetOrderItemsByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (os *OrderService) PlaceOrder(ctx context.Context, userEmail string, items []models.OrderItem) (orderID int, err error) {

	orderID, err = os.storage.PlaceOrder(ctx, userEmail, items)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return 0, service.ErrNotFound
		}
		return 0, err
	}
	return orderID, nil
}
