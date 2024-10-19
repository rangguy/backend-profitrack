package criteria

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"profitrack/middleware"
)

func Initiator(router *gin.Engine, db *gorm.DB) {
	repo := NewCriteriaRepository(db)
	service := NewCriteriaService(repo)

	api := router.Group("/api")
	api.Use(middleware.LoggingMiddleware())
	api.Use(middleware.JWTMiddleware())
	api.GET("/criterias", service.GetAllCriteriaService)
	api.GET("/criterias/:id", service.GetCriteriaByIdService)
	api.POST("/criterias", service.CreateCriteriaService)
	api.PUT("/criterias/:id", service.UpdateCriteriaService)
	api.DELETE("/criterias/:id", service.DeleteCriteriaService)
}
