package criteria

import (
	"time"
)

type Criteria struct {
	ID        int       `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Name      string    `gorm:"type:varchar(25);UNIQUE;not null" json:"name"`
	Weight    float64   `gorm:"double;not null" json:"weight"`
	Type      string    `gorm:"varchar(25);not null" json:"type"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type ResponseCriteria struct {
	ID     int     ` json:"id"`
	Name   string  `json:"name"`
	Weight float64 ` json:"weight"`
	Type   string  `json:"type"`
}
