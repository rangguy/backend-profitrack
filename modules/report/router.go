package report

import (
	"backend-profitrack/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Initiator(router *gin.Engine, db *gorm.DB) {
	repo := NewReportRepository(db)
	service := NewReportService(repo)

	api := router.Group("/api")
	api.Use(middleware.LoggingMiddleware())
	api.Use(middleware.JWTMiddleware())
	api.POST("/reports/:methodID", service.GetAllReportService)
	api.POST("/reports/export/:methodID", service.ExportExcelService)
	api.DELETE("/reports/:methodID", service.DeleteAllReportService)
}
