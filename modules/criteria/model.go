package criteria

import (
	"time"
)

type Criteria struct {
	ID         int       `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Name       string    `gorm:"type:varchar(255);UNIQUE;not null" json:"name"`
	Weight     float64   `gorm:"double;not null" json:"weight"`
	Type       string    `gorm:"varchar(255);not null" json:"type"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	ModifiedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"modified_at"`
}

type ResponseCriteria struct {
	ID     int     ` json:"id"`
	Name   string  `json:"name"`
	Weight float64 ` json:"weight"`
	Type   string  `json:"type"`
}
