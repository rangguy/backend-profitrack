package report

import (
	"backend-profitrack/modules/method"
	"backend-profitrack/modules/product"
	"time"
)

type Report struct {
	ID         int             `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	ProductID  int             `gorm:"integer;not null" json:"product_id"`
	MethodID   int             `gorm:"integer;not null" json:"method_id"`
	FinalScore float64         `gorm:"double" json:"final_score"`
	Period     string          `gorm:"varchar(255)" json:"period"`
	Product    product.Product `gorm:"foreignkey:ProductID" json:"product"`
	Method     method.Method   `gorm:"foreignkey:MethodID" json:"-"`
	CreatedAt  time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
