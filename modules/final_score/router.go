package final_score

import (
	"backend-profitrack/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Initiator(router *gin.Engine, db *gorm.DB) {
	repo := NewFinalScoreRepository(db)
	service := NewFinalScoreService(repo)

	api := router.Group("/api")
	api.Use(middleware.LoggingMiddleware())
	api.Use(middleware.JWTMiddleware())
	api.GET("/final_scores/:methodID", service.GetAllFinalScoreByMethodID)
}
