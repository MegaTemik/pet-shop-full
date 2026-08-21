package order

import (
	"context"
	"fmt"
	"go-pet-shop/internal/models"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

//go:generate go run github.com/vektra/mockery/v2 --name=Orders
type Orders interface {
	CreateOrder(ctx context.Context, order models.Order) (int, error)
	AddOrderItem(ctx context.Context, orderItem models.OrderItem) error
	GetOrderByID(ctx context.Context, id int) (models.Order, error)
	GetOrdersByUserEmail(ctx context.Context, email string) ([]models.Order, error)
	GetOrderItemsByOrderID(ctx context.Context, orderID int) ([]models.OrderItem, error)
	PlaceOrder(ctx context.Context, userEmail string, items []models.OrderItem) (orderID int, err error)
}

type PlaceOrderRequest struct {
	UserEmail string             `json:"user_email"`
	Items     []models.OrderItem `json:"items"`
}

type Handler struct {
	log     *slog.Logger
	storage Orders
}

func New(log *slog.Logger, storage Orders) *Handler {
	return &Handler{
		log:     log,
		storage: storage,
	}
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.orders.createOrder"
	log := h.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	log.Info("Creating new order", slog.String("url", r.URL.String()))

	var order models.Order
	if err := render.DecodeJSON(r.Body, &order); err != nil {
		log.Error("failed to decode request body", slog.Any("error", err))
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Invalid JSON payload",
		})
		return
	}

	if order.UserEmail == "" {
		log.Error("user email is empty")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "User email is required",
		})
		return
	}

	if order.TotalPrice < 0 {
		log.Error("total price is negative", slog.Float64("total_price", order.TotalPrice))
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Total price cannot be negative",
		})
		return
	}

	orderID, err := h.storage.CreateOrder(r.Context(), order)
	if err != nil {
		log.Error("failed to create order", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error":   "Internal server error",
			"message": "Failed to create order",
		})
		return
	}

	log.Info("Order created successfully", slog.String("user_email", order.UserEmail))

	order.ID = orderID
	render.JSON(w, r, map[string]interface{}{
		"message": "Order created successfully",
		"id":      orderID,
		"order":   order,
	})
}

func (h *Handler) AddOrderItem(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.orders.addOrderItem"

	log := h.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	log.Info("Adding order item", slog.String("url", r.URL.String()))

	orderIDstr := chi.URLParam(r, "id")
	orderID, err := strconv.Atoi(orderIDstr)
	if err != nil {
		log.Error("invalid order ID format", slog.Any("error", err), slog.String("id", orderIDstr))
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Order ID must be a number",
		})
		return
	}

	var orderItem models.OrderItem
	if err := render.DecodeJSON(r.Body, &orderItem); err != nil {
		log.Error("failed to decode request body", slog.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Invalid JSON payload",
		})
		return
	}

	if orderItem.Quantity < 0 {
		log.Error("quantity is negative", slog.Int("quantity", orderItem.Quantity))
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Quantity cannot be negative",
		})
		return
	}

	orderItem.OrderID = orderID

	err = h.storage.AddOrderItem(r.Context(), orderItem)
	if err != nil {
		log.Error("failed to add order item", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error":   "Internal server error",
			"message": "Failed to add order item",
		})
		return
	}

	log.Info("Order item added successfully",
		slog.Int("order_id", orderItem.OrderID),
		slog.Int("product_id", orderItem.ProductID),
	)

	render.JSON(w, r, map[string]interface{}{
		"message":    "Order item added successfully",
		"order_item": orderItem,
	})
}

func (h *Handler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.orders.getOrderByID"

	log := h.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	log.Info("Getting order by ID", slog.String("url", r.URL.String()))

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		log.Error("empty id")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "ID is required",
		})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error("invalid id format", slog.Any("error", err), slog.String("id", idStr))
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Order ID must be a number",
		})
		return
	}

	order, err := h.storage.GetOrderByID(r.Context(), id)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") ||
			strings.Contains(strings.ToLower(err.Error()), "no rows") ||
			strings.Contains(strings.ToLower(err.Error()), "rows affected: 0") {
			log.Warn("order not found", slog.Int("id", id))
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, map[string]string{
				"error":   "Not found",
				"message": fmt.Sprintf("Order with ID %d does not exist", id),
				"id":      idStr,
			})
			return
		}

		log.Error("failed to retrieve order", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error":   "Internal server error",
			"message": "Failed to retrieve order",
		})
		return
	}

	log.Info("Order retrieved successfully",
		slog.Int("id", id),
		slog.String("url", r.URL.String()),
	)

	render.JSON(w, r, map[string]interface{}{
		"status": "Order retrieved successfully",
		"id":     id,
		"order":  order,
	})
}

func (h *Handler) GetOrdersByUserEmail(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.orders.getOrdersByUserEmail"

	log := h.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	log.Info("Getting orders", slog.String("url", r.URL.String()))

	email := r.URL.Query().Get("email")
	if email == "" {
		log.Error("empty email")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Order email is required",
		})
		return
	}

	orders, err := h.storage.GetOrdersByUserEmail(r.Context(), email)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") ||
			strings.Contains(strings.ToLower(err.Error()), "no rows") ||
			strings.Contains(strings.ToLower(err.Error()), "rows affected: 0") {
			log.Warn("orders not found", slog.String("email", email))
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, map[string]interface{}{
				"error":   "Not found",
				"message": fmt.Sprintf("Orders with email %s does not exist", email),
				"email":   email,
			})
			return
		}
		log.Error("failed to get orders by email", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error":   "Internal server error",
			"message": "Failed to retrieve orders",
		})
	}

	log.Info("Retrieved orders successfully",
		slog.String("url", r.URL.String()),
		slog.String("email", email),
	)

	render.JSON(w, r, map[string]interface{}{
		"status": "Orders retrieved successfully",
		"email":  email,
		"orders": orders,
	})
}

func (h *Handler) GetOrderItemsByOrderID(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.orders.getOrderItemsByOrderID"

	log := h.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	log.Info("Getting order items", slog.String("url", r.URL.String()))

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		log.Error("empty order id in URL")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "ID is required",
		})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error("invalid id format", slog.Any("error", err), slog.String("id", idStr))
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "ID must be a number",
		})
		return
	}

	orderItems, err := h.storage.GetOrderItemsByOrderID(r.Context(), id)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") ||
			strings.Contains(strings.ToLower(err.Error()), "no rows") ||
			strings.Contains(strings.ToLower(err.Error()), "rows affected: 0") {
			log.Warn("order items not found", slog.Int("id", id))
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, map[string]interface{}{
				"error":   "Not found",
				"message": fmt.Sprintf("Order items with order ID %d do not exist", id),
				"id":      id,
			})
			return
		}

		log.Error("failed to get order items by order ID", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error":   "Internal server error",
			"message": "Failed to retrieve order items",
		})
		return
	}

	log.Info("Retrieved order items successfully",
		slog.String("url", r.URL.String()),
		slog.Int("id", id),
	)

	render.JSON(w, r, map[string]interface{}{
		"status": "Order items retrieved successfully",
		"id":     id,
		"items":  orderItems,
	})
}

func (h *Handler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.orders.PlaceOrder"

	log := h.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	log.Info("Placing new order", slog.String("url", r.URL.String()))

	var orderItems PlaceOrderRequest
	if err := render.DecodeJSON(r.Body, &orderItems); err != nil {
		log.Error("failed to decode request body", slog.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Invalid JSON payload",
		})
		return
	}

	for _, item := range orderItems.Items {
		if item.Quantity < 0 {
			log.Error("quantity is negative", slog.Int("quantity", item.Quantity))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, map[string]string{
				"error":   "Bad request",
				"message": "Quantity cannot be negative",
			})
			return
		}
	}

	orderID, err := h.storage.PlaceOrder(r.Context(), orderItems.UserEmail, orderItems.Items)
	if err != nil {
		log.Error("failed to place order", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error":   "Internal server error",
			"message": "Failed to place order",
		})
		return
	}

	log.Info("Order placed successfully", slog.String("user_email", orderItems.UserEmail), slog.Int("order_id", orderID))

	render.JSON(w, r, map[string]interface{}{
		"status":     "Order placed successfully",
		"order_id":   orderID,
		"user_email": orderItems.UserEmail,
	})
}
