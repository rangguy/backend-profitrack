package report

import (
	"backend-profitrack/modules/method"
	"backend-profitrack/modules/product"
	"time"
)

type Report struct {
	ID         int           `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	MethodID   int           `gorm:"integer;not null" json:"method_id"`
	ReportCode string        `gorm:"varchar(50);not null" json:"report_code"`
	TotalData  int           `gorm:"double" json:"total_data"`
	Method     method.Method `gorm:"foreignkey:MethodID" json:"-"`
	CreatedAt  time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type ReportDetail struct {
	ID         int             `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	MethodID   int             `gorm:"integer;not null" json:"method_id"`
	ProductID  int             `gorm:"integer;not null" json:"product_id"`
	ReportID   int             `gorm:"integer;not null" json:"report_id"`
	FinalScore float64         `gorm:"double" json:"final_scores"`
	Method     method.Method   `gorm:"foreignkey:MethodID" json:"-"`
	Product    product.Product `gorm:"foreignkey:ProductID" json:"product"`
	Report     Report          `gorm:"foreignkey:ReportID" json:"-"`
	CreatedAt  time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
