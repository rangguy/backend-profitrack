package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"os"
	"profitrack/database"
	"profitrack/modules/category"
	"profitrack/modules/criteria"
	"profitrack/modules/product"
	"profitrack/modules/user"
)

func main() {
	db := database.ConnectDatabase()

	InitiateRouter(db)
}

func InitiateRouter(db *gorm.DB) {
	router := gin.Default()

	user.Initiator(router, db)
	category.Initiator(router, db)
	product.Initiator(router, db)
	criteria.Initiator(router, db)

	err := router.Run(":" + os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
