package handlers

import (
	"e-commerce/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"bytes"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type mockProductService struct {
	products []models.Product
	err      error
}

func (m *mockProductService) GetProducts() ([]models.Product, error) {
	return m.products, m.err
}

func (m *mockProductService) CreateProduct(product *models.Product) error {
	if m.err != nil {
		return m.err
	}
	m.products = append(m.products, *product)
	return nil
}

func (m *mockProductService) UpdateProduct(product *models.Product) error {
	return m.err
}

func (m *mockProductService) DeleteProduct(id uint) error {
	return m.err
}

func TestGetProductsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockProductService{
		products: []models.Product{
			{ID: 1, Name: "Test", Price: 10},
		},
	}

	handler := NewProductHandler(mockService)

	r := gin.Default()
	r.GET("/products", handler.GetProducts)

	req, _ := http.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestCreateProductHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockProductService{}

	handler := NewProductHandler(mockService)

	r := gin.Default()
	r.POST("/products", handler.CreateProduct)

	body, _ := json.Marshal(models.Product{
		Name:  "Test",
		Price: 10,
	})

	req, _ := http.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
}

func TestUpdateProductHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockProductService{}

	handler := NewProductHandler(mockService)

	r := gin.Default()
	r.PUT("/products/:id", handler.UpdateProduct)

	body, _ := json.Marshal(models.Product{
		Name:  "Updated",
		Price: 50,
	})

	req, _ := http.NewRequest(http.MethodPut, "/products/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestDeleteProductHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockProductService{}

	handler := NewProductHandler(mockService)

	r := gin.Default()
	r.DELETE("/products/:id", handler.DeleteProduct)

	req, _ := http.NewRequest(http.MethodDelete, "/products/1", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}
