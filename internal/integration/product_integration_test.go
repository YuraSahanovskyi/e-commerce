package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"e-commerce/internal/handlers"
	"e-commerce/internal/models"
	"e-commerce/internal/repositories"
	"e-commerce/internal/services"

	"github.com/gin-gonic/gin"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, func()) {
	ctx := context.Background()

	container, err := tcpostgres.RunContainer(ctx,
		tcpostgres.WithDatabase("testdb"),
		tcpostgres.WithUsername("user"),
		tcpostgres.WithPassword("pass"),
	)
	if err != nil {
		t.Fatalf("failed to start container: %v", err)
	}

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	var db *gorm.DB

	for i := 0; i < 10; i++ {
		db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}

	if err != nil {
		t.Fatalf("failed to connect to db: %v", err)
	}

	// міграції
	if err := db.AutoMigrate(&models.Product{}); err != nil {
		t.Fatalf("migration failed: %v", err)
	}

	cleanup := func() {
		_ = container.Terminate(ctx)
	}

	return db, cleanup
}

func TestCreateProduct_Integration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewProductRepository(db)
	service := services.NewProductService(repo)
	handler := handlers.NewProductHandler(service)

	r := gin.Default()
	r.POST("/products", handler.CreateProduct)
	r.GET("/products", handler.GetProducts)

	body, _ := json.Marshal(models.Product{
		Name:  "Integration Test",
		Price: 100,
	})

	req, _ := http.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	req2, _ := http.NewRequest(http.MethodGet, "/products", nil)
	w2 := httptest.NewRecorder()

	r.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w2.Code)
	}

	var products []models.Product
	err := json.Unmarshal(w2.Body.Bytes(), &products)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(products) != 1 {
		t.Fatalf("expected 1 product, got %d", len(products))
	}

	if products[0].Name != "Integration Test" || products[0].Price != 100 {
		t.Fatalf("product fields mismatch: %+v", products[0])
	}
}
