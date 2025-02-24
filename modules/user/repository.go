package user

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"time"
)

type Repository interface {
	LoginRepository(username string) (User, error)
	GetUserByIDRepository(userID int) (User, error)
	UpdateByIDRepository(user *User) (err error)
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) Repository {
	var count int64
	db.Model(&User{}).Where("username = ?", "admin").Count(&count)

	if count == 0 {
		var password []byte
		password, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("failed to hash password: ", err)
		}
		adminUser := User{
			Username:  "admin",
			Password:  string(password),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		db.Create(&adminUser)
		log.Println("Admin user created.")
	} else {
		log.Println("Admin user already exists.")
	}

	return &userRepository{
		DB: db,
	}
}

func (r *userRepository) LoginRepository(username string) (User, error) {
	var user User
	err := r.DB.Where("username = ?", username).First(&user).Error
	return user, err
}

func (r *userRepository) GetUserByIDRepository(userID int) (User, error) {
	var user User
	err := r.DB.Where("id = ?", userID).First(&user).Error
	return user, err
}

func (r *userRepository) UpdateByIDRepository(user *User) (err error) {
	err = r.DB.Save(user).Error
	return err
}
