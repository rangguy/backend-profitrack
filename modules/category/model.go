package category

import "time"

type Category struct {
	ID         int       `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Name       string    `gorm:"varchar(255);UNIQUE" json:"name"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	ModifiedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"modified_at"`
}
