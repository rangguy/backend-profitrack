package product

import (
	"backend-profitrack/modules/category"
	"time"
)

type Product struct {
	ID           int               `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Name         string            `gorm:"varchar(255);UNIQUE;not null" json:"name"`
	PurchaseCost int               `gorm:"integer;not null" json:"purchase_cost"`
	PriceSale    int               `gorm:"integer;not null" json:"price_sale"`
	Profit       int               `gorm:"integer" json:"profit"`
	Unit         string            `gorm:"varchar(255)" json:"unit"`
	Stock        int               `gorm:"integer" json:"stock"`
	Sold         int               `gorm:"integer" json:"sold"`
	CategoryID   int               `gorm:"integer" json:"category_id"`
	Category     category.Category `gorm:"foreignkey:CategoryID" json:"category"`
	CreatedAt    time.Time         `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time         `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type ResponseProduct struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	PurchaseCost int    `json:"purchase_cost"`
	PriceSale    int    `json:"price_sale"`
	Profit       int    `json:"profit"`
	Unit         string `json:"unit"`
	Stock        int    `json:"stock"`
	Sold         int    `json:"sold"`
	CategoryID   int    `json:"category_id"`
	CategoryName string `json:"category_name"`
}

type ExcelProduct struct {
	Name         string `validate:"required"`
	PurchaseCost int    `validate:"required"`
	PriceSale    int    `validate:"required"`
	Stock        int    `validate:"required"`
	Sold         int    `validate:"required"`
	CategoryName string `validate:"required"`
}
