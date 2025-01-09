package main

import (
	"backend-profitrack/database"
	"backend-profitrack/database/migrations"
	"backend-profitrack/middleware"
	"backend-profitrack/modules/category"
	"backend-profitrack/modules/criteria"
	"backend-profitrack/modules/product"
	"backend-profitrack/modules/score_smart"
	"backend-profitrack/modules/user"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"os"
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
	score_smart.Initiator(router, db)

	err := router.Run(":" + os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
