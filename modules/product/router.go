package product

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Initiator(router *gin.Engine, db *gorm.DB) {
	repo := NewProductRepository(db)
	service := NewProductService(repo)

	api := router.Group("/api")
	api.GET("/products", service.GetAllProductService)
	api.GET("/products/:id", service.GetProductByIdService)
	api.POST("/products", service.CreateProductService)
	api.PUT("/products/:id", service.UpdateProductService)
	api.DELETE("/products/:id", service.DeleteProductService)
}
