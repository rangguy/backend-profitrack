package database

import (
	"fmt"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"profitrack/modules/category"
	"profitrack/modules/user"
	"time"
)

var DBConnection *gorm.DB

func ConnectDatabase() {
	err := godotenv.Load("config/.env")
	if err != nil {
		panic("Error loading .env file")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("PGHOST"),
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
		os.Getenv("PGDATABASE"),
		os.Getenv("PGPORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("Successfully connected to database")

	err = db.AutoMigrate(&user.User{}, &category.Category{})
	if err != nil {
		panic(err)
	}

	var count int64
	db.Model(&user.User{}).Where("username = ?", "admin").Count(&count)

	if count == 0 {
		var password []byte
		password, err = bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("failed to hash password: ", err)
		}
		adminUser := user.User{
			Username:   "admin",
			Password:   string(password),
			CreatedAt:  time.Now(),
			ModifiedAt: time.Now(),
		}
		db.Create(&adminUser)
		log.Println("Admin user created.")
	} else {
		log.Println("Admin user already exists.")
	}

	DBConnection = db
	fmt.Println("Migration success")
}
