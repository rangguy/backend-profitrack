package user

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"profitrack/config"
	"profitrack/helpers"
	"time"
)

type Service interface {
	LoginService(ctx *gin.Context)
	LogoutService(ctx *gin.Context)
}

type userService struct {
	repository Repository
}

func NewUserService(repository Repository) Service {
	return &userService{
		repository,
	}
}

func (service *userService) LoginService(ctx *gin.Context) {
	var userRequest LoginRequest

	if err := ctx.ShouldBind(&userRequest); err != nil {
		response := map[string]string{"error": "failed to read body"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	user, err := service.repository.LoginRepository(userRequest.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := map[string]string{"message": "Username atau password salah"}
			helpers.ResponseJSON(ctx, http.StatusUnauthorized, response)
			return
		}
		response := map[string]string{"message": err.Error()}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userRequest.Password)); err != nil {
		response := map[string]string{"message": "Username atau password salah"}
		helpers.ResponseJSON(ctx, http.StatusUnauthorized, response)
		return
	}

	expiredTime := time.Now().Add(time.Hour * 2)
	claims := &config.JWTClaim{
		Username: userRequest.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "profitrack",
			ExpiresAt: jwt.NewNumericDate(expiredTime),
		},
	}

	tokenAlgo := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenAlgo.SignedString(config.JWT_KEY)
	if err != nil {
		response := map[string]string{"message": err.Error()}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    token,
		HttpOnly: false,
	})

	response := map[string]string{"message": "login berhasil"}
	helpers.ResponseJSON(ctx, http.StatusOK, response)
}

func (service *userService) LogoutService(ctx *gin.Context) {
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
	})

	response := map[string]string{"message": "logout berhasil"}
	helpers.ResponseJSON(ctx, http.StatusOK, response)
}
