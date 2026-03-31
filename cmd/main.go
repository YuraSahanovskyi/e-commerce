package main

import (
	"context"
	"e-commerce/internal/config"
	"e-commerce/internal/db"
	"e-commerce/internal/handlers"
	"e-commerce/internal/logger"
	"e-commerce/internal/repositories"
	"e-commerce/internal/services"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	logger.Init()

	if err := godotenv.Load(); err != nil {
		logger.Log.Info(".env not found, using system env")
	} else {
		logger.Log.Info(".env loaded")
	}

	db.InitDB()
	db.RunMigrations()

	r := gin.New()
	r.Use(gin.Recovery())

	r.Use(func(c *gin.Context) {
		c.Next()

		msg := fmt.Sprintf("%s %s %d",
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
		)

		logger.Log.Info(msg)
	})

	productRepo := repositories.NewProductRepository(db.DB)

	productService := services.NewProductService(productRepo)

	productHandler := handlers.NewProductHandler(productService)

	r.GET("/health", handlers.Health)
	r.GET("/products", productHandler.GetProducts)
	r.POST("/products", productHandler.CreateProduct)
	r.PUT("/products/:id", productHandler.UpdateProduct)
	r.DELETE("/products/:id", productHandler.DeleteProduct)

	port := config.GetEnv("PORT")

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		logger.Log.Info("Starting server on port " + port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Error("listen:" + err.Error())
		}
	}()

	// graceful shutdown
	gracefulShutdown(srv)
}

func gracefulShutdown(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	if sig == syscall.SIGTERM {
		logger.Log.Info("SIGTERM received. Starting graceful shutdown...")
	} else {
		logger.Log.Info("Signal received. Starting shutdown...")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Error("Server forced to shutdown: " + err.Error())
	}

	sqlDB, err := db.DB.DB()
	if err == nil {
		sqlDB.Close()
		logger.Log.Info("Database connection closed")
	}

	logger.Log.Info("Server exiting")
	os.Exit(0)
}
