package product

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"profitrack/middleware"
	"profitrack/modules/category"
)

func Initiator(router *gin.Engine, db *gorm.DB) {
	repo := NewProductRepository(db)
	categoryRepo := category.NewCategoryRepository(db)
	service := NewProductService(repo, categoryRepo)

	api := router.Group("/api")
	api.Use(middleware.LoggingMiddleware())
	api.Use(middleware.JWTMiddleware())
	api.GET("/products", service.GetAllProductService)
	api.GET("/products/:id", service.GetProductByIdService)
	api.POST("/products", service.CreateProductService)
	api.PUT("/products/:id", service.UpdateProductService)
	api.DELETE("/products/:id", service.DeleteProductService)
	api.POST("/products/import", service.ImportExcelService)
	api.GET("/products/export", service.ExportExcelService)
}
