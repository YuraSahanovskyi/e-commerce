package models

type Product struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Description string `json:"description"`
}