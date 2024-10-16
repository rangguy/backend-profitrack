package user

import (
	"gorm.io/gorm"
)

type Repository interface {
	LoginRepository(username string) (User, error)
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) Repository {
	return &userRepository{
		DB: db,
	}
}

func (r *userRepository) LoginRepository(username string) (User, error) {
	var user User
	err := r.DB.Where("username = ?", username).First(&user).Error
	return user, err
}
