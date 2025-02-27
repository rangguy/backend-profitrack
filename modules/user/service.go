package user

import (
	"backend-profitrack/config"
	"backend-profitrack/helpers"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

type Service interface {
	LoginService(ctx *gin.Context)
	LogoutService(ctx *gin.Context)
	UpdatePasswordService(ctx *gin.Context)
	CountUserService(ctx *gin.Context)
}

type userService struct {
	repository Repository
}

func NewUserService(repository Repository) Service {
	return &userService{
		repository,
	}
}

func (service *userService) CountUserService(ctx *gin.Context) {
	result, err := service.repository.CountUserRepository()
	if err != nil {
		response := map[string]string{"error": "gagal menghitung jumlah data user"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	response := map[string]int{"count": int(result)}
	helpers.ResponseJSON(ctx, http.StatusOK, response)
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
		UserID:   user.ID,
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
		Value:    token,
		Path:     "/",
		Expires:  expiredTime,
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	response := map[string]string{"message": "login berhasil"}
	helpers.ResponseJSON(ctx, http.StatusOK, response)
}

func (service *userService) LogoutService(ctx *gin.Context) {
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    "",
		HttpOnly: false,
		MaxAge:   -1,
	})

	response := map[string]string{"message": "logout berhasil"}
	helpers.ResponseJSON(ctx, http.StatusOK, response)
}

func (service *userService) UpdatePasswordService(ctx *gin.Context) {
	cookie, err := ctx.Cookie("token")
	if err != nil {
		response := map[string]string{"error": "unauthorized"}
		helpers.ResponseJSON(ctx, http.StatusUnauthorized, response)
		return
	}

	claims := &config.JWTClaim{}
	token, err := jwt.ParseWithClaims(cookie, claims, func(token *jwt.Token) (interface{}, error) {
		return config.JWT_KEY, nil
	})
	if err != nil || !token.Valid {
		response := map[string]string{"error": "invalid token"}
		helpers.ResponseJSON(ctx, http.StatusUnauthorized, response)
		return
	}

	userID := claims.UserID

	var req UpdatePasswordRequest

	// Bind JSON body ke struct
	if err = ctx.ShouldBindJSON(&req); err != nil {
		response := map[string]string{"error": "failed to read body"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	// Ambil user berdasarkan ID dari token
	user, err := service.repository.GetUserByIDRepository(userID)
	if err != nil {
		response := map[string]string{"error": "user not found"}
		helpers.ResponseJSON(ctx, http.StatusNotFound, response)
		return
	}

	// Cek apakah old password sesuai
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		response := map[string]string{"error": "password lama tidak sesuai, silahkan coba lagi"}
		helpers.ResponseJSON(ctx, http.StatusUnauthorized, response)
		return
	}

	// Hash password baru
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Println("failed to hash password:", err)
		response := map[string]string{"error": "failed to hash password"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	if req.OldPassword == req.NewPassword {
		response := map[string]string{"error": "Password baru tidak boleh sama dengan password lama, silahkan coba lagi"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	user.Password = string(hashedPassword)

	// Update password di database
	if err = service.repository.UpdateByIDRepository(&user); err != nil {
		response := map[string]string{"error": "gagal mengubah password, silahkan coba lagi"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	// Berhasil
	response := map[string]string{"message": "password berhasil diubah"}
	helpers.ResponseJSON(ctx, http.StatusOK, response)
}
