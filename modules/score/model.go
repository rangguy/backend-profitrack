package score

import (
	"backend-profitrack/modules/criteria"
	"backend-profitrack/modules/method"
	"backend-profitrack/modules/product"
	"time"
)

type Score struct {
	ID             int               `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	ProductID      int               `gorm:"integer;not null" json:"product_id"`
	CriteriaID     int               `gorm:"integer;not null" json:"criteria_id"`
	MethodID       int               `gorm:"integer;not null" json:"method_id"`
	Score          float64           `gorm:"double" json:"score"`
	NormalizeScore float64           `gorm:"double" json:"normalize_score"`
	ScoreOne       float64           `gorm:"double" json:"score_one"`
	ScoreTwo       float64           `gorm:"double" json:"score_two"`
	Product        product.Product   `gorm:"foreignkey:ProductID" json:"-"`
	Criteria       criteria.Criteria `gorm:"foreignkey:CriteriaID" json:"-"`
	Method         method.Method     `gorm:"foreignkey:MethodID" json:"-"`
	CreatedAt      time.Time         `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time         `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
