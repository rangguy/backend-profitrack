package value

import (
	"backend-profitrack/middleware"
	"backend-profitrack/modules/criteria"
	"backend-profitrack/modules/product"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Initiator(router *gin.Engine, db *gorm.DB) {
	repo := NewValueRepository(db)
	productRepo := product.NewProductRepository(db)
	criteriaRepo := criteria.NewCriteriaRepository(db)
	service := NewValueService(repo, productRepo, criteriaRepo)

	api := router.Group("/api")
	api.Use(middleware.LoggingMiddleware())
	api.Use(middleware.JWTMiddleware())
	api.GET("/values", service.GetAllValueService)
	api.POST("/values", service.CreateValueService)
	api.DELETE("/values", service.DeleteAllValuesService)
}
