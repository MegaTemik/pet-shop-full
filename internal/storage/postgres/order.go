package postgres

import (
	"context"
	"errors"
	"fmt"
	"go-pet-shop/internal/models"
	"go-pet-shop/internal/storage"

	"github.com/jackc/pgx/v5"
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
	const fn = "storage.postgres.order.AddOrderItem"

	_, err := s.db.Exec(ctx,
		`INSERT INTO order_items (order_id, product_id, quantity) VALUES($1, $2, $3)`,
		orderItem.OrderID, orderItem.ProductID, orderItem.Quantity)

	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	return nil
}

func (s *Storage) GetOrderByID(ctx context.Context, id int) (models.Order, error) {
	const fn = "storage.postgres.order.GetOrderByID"

	var order models.Order
	err := s.db.QueryRow(ctx,
		`SELECT id, user_email, total_price, created_at FROM orders WHERE id = $1`, id).
		Scan(&order.ID, &order.UserEmail, &order.TotalPrice, &order.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Order{}, fmt.Errorf("%s: %w", fn, storage.ErrNotFound)
		}
		return models.Order{}, fmt.Errorf("%s: %w", fn, err)
	}
	return order, nil

}

func (s *Storage) GetOrdersByUserEmail(ctx context.Context, email string) ([]models.Order, error) {
	const fn = "storage.postgres.order.GetOrdersByUserEmail"

	rows, err := s.db.Query(ctx,
		`SELECT o.id, o.user_email, o.total_price, o.created_at FROM orders o JOIN users u ON u.email = o.user_email WHERE u.email = $1`, email)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.UserEmail, &o.TotalPrice, &o.CreatedAt); err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		orders = append(orders, o)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return orders, nil
}

func (s *Storage) GetOrderItemsByOrderID(ctx context.Context, orderID int) ([]models.OrderItem, error) {
	const fn = "storage.postgres.order.GetOrderItemsByOrderID"

	rows, err := s.db.Query(ctx,
		`SELECT id, order_id, product_id, quantity FROM order_items WHERE order_id = $1`, orderID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	defer rows.Close()

	orderItems := make([]models.OrderItem, 0)
	for rows.Next() {
		var oi models.OrderItem
		if err := rows.Scan(&oi.ID, &oi.OrderID, &oi.ProductID, &oi.Quantity); err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		orderItems = append(orderItems, oi)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return orderItems, nil
}

func (s *Storage) PlaceOrder(ctx context.Context, userEmail string, items []models.OrderItem) (orderID int, err error) {
	const fn = "storage.postgres.order.PlaceOrder"

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var totalSum float64
	for _, item := range items {
		var price float64
		err := tx.QueryRow(ctx,
			`UPDATE products SET stock = stock - $1 WHERE id = $2 AND stock >= $1 RETURNING price`,
			item.Quantity, item.ProductID).Scan(&price)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return 0, fmt.Errorf("%s: %w", fn, storage.ErrNotFound)
			}
			return 0, fmt.Errorf("%s: %w", fn, err)
		}
		totalSum += price * float64(item.Quantity)
	}

	err = tx.QueryRow(ctx,
		`INSERT INTO orders (user_email, total_price) VALUES ($1, $2) RETURNING id`, userEmail, totalSum).Scan(&orderID)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	for _, item := range items {
		_, err := tx.Exec(ctx,
			`INSERT INTO order_items (order_id, product_id, quantity) VALUES ($1, $2, $3)`, orderID, item.ProductID, item.Quantity)
		if err != nil {
			return 0, fmt.Errorf("%s: %w", fn, err)
		}
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO transactions (order_id, amount, status) VALUES ($1, $2, $3)`, orderID, totalSum, "success")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	return
}

func (s *Storage) GetUserOrderHistory(ctx context.Context, email string) ([]models.OrderDetail, error) {
	const fn = "storage.postgres.order.GetUserOrderHistory"

	rows, err := s.db.Query(ctx,
		`SELECT o.id, o.user_email, o.total_price, o.created_at, t.status
		FROM orders o
		JOIN users u ON o.user_email = u.email
		JOIN transactions t ON o.id = t.order_id
		WHERE u.email = $1`, email)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	defer rows.Close()

	var ordersHistory []models.OrderDetail
	for rows.Next() {
		var order models.Order
		var status string
		if err := rows.Scan(&order.ID, &order.UserEmail, &order.TotalPrice, &order.CreatedAt, &status); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []models.OrderDetail{}, fmt.Errorf("%s: %w", fn, storage.ErrNotFound)
			}
			return nil, fmt.Errorf("%s: %w", fn, err)
		}

		orderItems, err := s.GetOrderItemsByOrderID(ctx, order.ID)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}

		ordersDetail := models.OrderDetail{
			Order:  order,
			Items:  orderItems,
			Status: status,
		}
		ordersHistory = append(ordersHistory, ordersDetail)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	return ordersHistory, nil
}
