package score_smart

import (
	"backend-profitrack/modules/criteria"
	"backend-profitrack/modules/product"
	"time"
)

type ScoreSmart struct {
	ID         int               `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	ProductID  int               `gorm:"integer;not null" json:"product_id"`
	CriteriaID int               `gorm:"integer;not null" json:"criteria_id"`
	Score      float64           `gorm:"double;not null" json:"score"`
	Period     time.Time         `gorm:"TIMESTAMP" json:"period"`
	Product    product.Product   `gorm:"foreignkey:ProductID" json:"-"`
	Criteria   criteria.Criteria `gorm:"foreignkey:CriteriaID" json:"-"`
	CreatedAt  time.Time         `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	ModifiedAt time.Time         `gorm:"default:CURRENT_TIMESTAMP" json:"modified_at"`
}
