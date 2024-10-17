package product

import (
	"profitrack/modules/category"
	"time"
)

type Product struct {
	ID           int               `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Name         string            `gorm:"varchar(255);UNIQUE;not null" json:"name"`
	NetProfit    int               `gorm:"integer;not null" json:"net_profit"`
	GrossProfit  int               `gorm:"integer;not null" json:"gross_profit"`
	GrossSale    int               `gorm:"integer;not null" json:"gross_sale"`
	PurchaseCost int               `gorm:"integer;not null" json:"purchase_cost"`
	InitialStock int               `gorm:"integer;not null" json:"initial_stock"`
	FinalStock   int               `gorm:"integer;not null" json:"final_stock"`
	CategoryID   int               `gorm:"integer;not null" json:"category_id"`
	Category     category.Category `gorm:"foreignkey:CategoryID" json:"-"`
	CreatedAt    time.Time         `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	ModifiedAt   time.Time         `gorm:"default:CURRENT_TIMESTAMP" json:"modified_at"`
}

type ResponseProduct struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	NetProfit    int    `json:"net_profit"`
	GrossProfit  int    `json:"gross_profit"`
	GrossSale    int    `json:"gross_sale"`
	PurchaseCost int    `json:"purchase_cost"`
	InitialStock int    `json:"initial_stock"`
	FinalStock   int    `json:"final_stock"`
	CategoryID   int    `json:"category_id"`
}
