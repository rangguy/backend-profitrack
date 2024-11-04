package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"os"
	"profitrack/database"
	"profitrack/database/migrations"
	"profitrack/middleware"
	"profitrack/modules/category"
	"profitrack/modules/criteria"
	"profitrack/modules/product"
	"profitrack/modules/user"
	"profitrack/modules/value"
)

func main() {
	db := database.ConnectDatabase()

	migrations.Migrations(db)
	InitiateRouter(db)
}

func InitiateRouter(db *gorm.DB) {
	router := gin.Default()

	router.Use(middleware.CORSMiddleware())
	user.Initiator(router, db)
	category.Initiator(router, db)
	product.Initiator(router, db)
	criteria.Initiator(router, db)
	value.Initiator(router, db)

	err := router.Run(":" + os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
