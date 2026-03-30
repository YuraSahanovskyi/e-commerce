package services

import (
	"e-commerce/internal/models"
	"testing"
)

type mockProductRepository struct {
	products []models.Product
	err      error
}

func (m *mockProductRepository) GetAll() ([]models.Product, error) {
	return m.products, m.err
}

func (m *mockProductRepository) Create(product *models.Product) error {
	if m.err != nil {
		return m.err
	}
	m.products = append(m.products, *product)
	return nil
}

func (m *mockProductRepository) Update(product *models.Product) error {
	return m.err
}

func (m *mockProductRepository) Delete(id uint) error {
	return m.err
}

func TestGetProducts(t *testing.T) {
	mockRepo := &mockProductRepository{
		products: []models.Product{
			{ID: 1, Name: "Test", Price: 10},
		},
	}

	service := NewProductService(mockRepo)

	products, err := service.GetProducts()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(products) != 1 {
		t.Fatalf("expected 1 product, got %d", len(products))
	}
}

func TestCreateProduct(t *testing.T) {
	mockRepo := &mockProductRepository{}

	service := NewProductService(mockRepo)

	product := &models.Product{
		Name:  "New",
		Price: 20,
	}

	err := service.CreateProduct(product)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(mockRepo.products) != 1 {
		t.Fatalf("product was not added")
	}
}

func TestUpdateProduct(t *testing.T) {
	mockRepo := &mockProductRepository{}

	service := NewProductService(mockRepo)

	product := &models.Product{
		ID:    1,
		Name:  "Updated",
		Price: 50,
	}

	err := service.UpdateProduct(product)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDeleteProduct(t *testing.T) {
	mockRepo := &mockProductRepository{}

	service := NewProductService(mockRepo)

	err := service.DeleteProduct(1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
