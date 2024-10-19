package value

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"profitrack/modules/criteria"
	"profitrack/modules/product"
)

func Initiator(router *gin.Engine, db *gorm.DB) {
	repo := NewValueRepository(db)
	productRepo := product.NewProductRepository(db)
	criteriaRepo := criteria.NewCriteriaRepository(db)
	service := NewValueService(repo, productRepo, criteriaRepo)

	api := router.Group("/api")
	api.GET("/values", service.GetAllValueService)
	api.POST("/values", service.CreateValueService)
	api.DELETE("/values", service.DeleteAllValuesService)
}
