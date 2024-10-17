package category

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Initiator(router *gin.Engine, db *gorm.DB) {
	repo := NewCategoryRepository(db)
	service := NewCategoryService(repo)

	api := router.Group("/api")
	api.GET("/categories", service.GetAllCategoryService)
	api.GET("/categories/:id", service.GetCategoryByIdService)
	api.POST("/categories", service.CreateCategoryService)
	api.PUT("/categories/:id", service.UpdateCategoryService)
	api.DELETE("/categories/:id", service.DeleteCategoryService)
}
