package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"os"
	"profitrack/database"
	"profitrack/modules/user"
)

func main() {
	database.ConnectDatabase()

	InitiateRouter(database.DBConnection)
}

func InitiateRouter(db *gorm.DB) {
	router := gin.Default()

	user.Initiator(router, db)

	err := router.Run(":" + os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
