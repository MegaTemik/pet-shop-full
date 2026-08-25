package main

import (
	"context"
	"go-pet-shop/internal/config"
	"go-pet-shop/internal/handlers"
	order_handler "go-pet-shop/internal/handlers/order"
	product_handler "go-pet-shop/internal/handlers/product"
	user_handler "go-pet-shop/internal/handlers/user"
	"go-pet-shop/internal/lib/logger"
	order_service "go-pet-shop/internal/service/order"
	product_service "go-pet-shop/internal/service/product"
	user_service "go-pet-shop/internal/service/user"
	"go-pet-shop/internal/storage/postgres"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	cfg := config.MustLoad()

	// Settings logger
	log := logger.SetupLogger(cfg.Env)
	log.Info("starting the project...", slog.String("env", cfg.Env))

	// Settings and started database
	ctx, cancel := context.WithTimeout(context.Background(), cfg.HTTPServer.DatabaseTimeout)
	defer cancel()

	storage, err := postgres.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Error("failed to init storage", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer func() {
		if err := storage.Close(); err != nil {
			log.Error("failed to close storage", slog.String("error", err.Error()))
		}
		log.Info("storage closed")
	}()

	productService := product_service.NewProductService(storage)
	userService := user_service.NewUserService(storage)
	orderService := order_service.NewOrderService(storage)

	// Init router
	router := chi.NewRouter()

	// Middlewares
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	//router.Use(middleware.URLFormat)
	router.Use(logger.CustomLogger(log))

	// Handlers
	productHandler := product_handler.New(log, productService)
	userHandler := user_handler.New(log, userService)
	orderHandler := order_handler.New(log, orderService)

	router.Get("/status", handlers.StatusHandler)

	router.Post("/users", userHandler.CreateUser)
	router.Get("/users", userHandler.GetAllUsers)
	router.Get("/users/{email}", userHandler.GetUserByEmail)

	router.Post("/products", productHandler.CreateProduct)
	router.Get("/products", productHandler.GetAllProducts)
	router.Get("/products/{id}", productHandler.GetProductByID)
	router.Put("/products/{id}", productHandler.UpdateProduct)
	router.Delete("/products/{id}", productHandler.DeleteProduct)

	router.Post("/orders", orderHandler.CreateOrder)
	router.Post("/orders/{id}/items", orderHandler.AddOrderItem)
	router.Get("/orders/{id}", orderHandler.GetOrderByID)
	router.Get("/users/orders", orderHandler.GetOrdersByUserEmail)

	router.Post("/checkout", orderHandler.PlaceOrder)

	// Settings and started server
	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	// Запуск сервера в горутине для graceful shutdown
	go func() {
		log.Info("Starting server on", slog.String("address", cfg.HTTPServer.Address))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Server error", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop // Ждем сигнал завершения

	log.Info("Shutting down server...")

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server shutdown error", slog.String("error", err.Error()))
	}

	log.Info("Server stopped gracefully")
}
