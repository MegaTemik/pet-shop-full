package postgres

import (
	"context"
	"errors"
	"fmt"
	"go-pet-shop/internal/models"
)

var (
	ErrOrderNotFound = errors.New("order not found")
)

func (s *Storage) CreateOrder(ctx context.Context, order models.Order) (int, error) {
	const fn = "storage.postgres.order.CreateOrder"

	var orderID int
	err := s.db.QueryRow(ctx,
		`INSERT INTO orders (user_email, total_price) VALUES ($1, $2) RETURNING id`,
		order.UserEmail, order.TotalPrice).Scan(&orderID)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	return orderID, nil
}

func (s *Storage) AddOrderItem(ctx context.Context, orderItem models.OrderItem) error {
	const fn = "storage.postgres.order.AddorderItem"

	_, err := s.db.Exec(ctx,
		`INSERT INTO order_items (order_id, product_id, quantity) VALUES($1, $2, $3)`,
		orderItem.OrderID, orderItem.ProductID, orderItem.Quantity)

	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	return nil
}

// CREATE TABLE orders (
//     id SERIAL PRIMARY KEY,
//     user_email TEXT REFERENCES users(email),
//     total_price NUMERIC NOT NULL,
//     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
// );

// CREATE TABLE order_items (
//     id SERIAL PRIMARY KEY,
//     order_id INT REFERENCES orders(id),
//     product_id INT REFERENCES products(id),
//     quantity INT NOT NULL
// );
