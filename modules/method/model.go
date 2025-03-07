package method

import "time"

type Method struct {
	ID        int       `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Name      string    `gorm:"varchar(25)" json:"name"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type ResponseMethod struct {
	ID   int    ` json:"id"`
	Name string `json:"name"`
}
