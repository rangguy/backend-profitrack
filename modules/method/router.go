package method

import (
	"backend-profitrack/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Initiator(router *gin.Engine, db *gorm.DB) {
	repo := NewMethodRepository(db)
	service := NewMethodService(repo)

	api := router.Group("/api")
	api.Use(middleware.LoggingMiddleware())
	api.Use(middleware.JWTMiddleware())
	api.GET("/methods", service.GetAllMethodService)
	api.GET("/methods/:id", service.GetMethodByIdService)
	api.DELETE("/methods/:id", service.DeleteMethodService)
}
