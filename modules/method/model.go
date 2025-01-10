package method

import "time"

type Method struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ResponseMethod struct {
	ID   int    ` json:"id"`
	Name string `json:"name"`
}