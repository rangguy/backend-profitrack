package criteria_score

import (
	"backend-profitrack/modules/criteria"
	"backend-profitrack/modules/product"
	"time"
)

type CriteriaScore struct {
	ID         int               `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	ProductID  int               `gorm:"smallint;not null" json:"product_id"`
	CriteriaID int               `gorm:"smallint;not null" json:"criteria_id"`
	Score      float64           `gorm:"double" json:"score"`
	Product    product.Product   `gorm:"foreignkey:ProductID" json:"-"`
	Criteria   criteria.Criteria `gorm:"foreignkey:CriteriaID" json:"-"`
	CreatedAt  time.Time         `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time         `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
