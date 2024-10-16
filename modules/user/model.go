package user

import "time"

type User struct {
	ID         int       `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Username   string    `gorm:"varchar(255);UNIQUE" json:"username"`
	Password   string    `gorm:"varchar(255)" json:"password"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	ModifiedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"modified_at"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}