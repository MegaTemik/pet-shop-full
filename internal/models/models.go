package models

import "time"

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"` // количество на складе
}

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Order struct {
	ID         int       `json:"id"`
	UserEmail  string    `json:"user_email"`
	TotalPrice float64   `json:"total_price"`
	CreatedAt  time.Time `json:"created_at"`
}

type OrderItem struct {
	ID        int `json:"id"`
	OrderID   int `json:"order_id"`
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type OrderDetail struct {
	Order      Order       `json:"order"`
	OrderItems []OrderItem `json:"order_items"`
}
