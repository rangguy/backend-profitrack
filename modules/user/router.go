package user

import (
	"backend-profitrack/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Initiator(router *gin.Engine, db *gorm.DB) {
	repo := NewUserRepository(db)
	service := NewUserService(repo)

	api := router.Group("/api")
	api.POST("/login", service.LoginService)
	api.GET("/logout", service.LogoutService)

	api.Use(middleware.LoggingMiddleware())
	api.Use(middleware.JWTMiddleware())
	api.PUT("/user", service.UpdatePasswordService)
	api.GET("/user/count", service.CountUserService)
}
