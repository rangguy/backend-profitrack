package middleware

import (
	"backend-profitrack/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

var mySigningKey = config.JWT_KEY

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil token dari Authorization Header
		authHeader := c.Request.Header.Get("Authorization")
		var tokenString string

		if authHeader != "" {
			tokenString = strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		} else {
			// Jika tidak ada di header, coba baca dari cookie
			cookie, err := c.Cookie("token")
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is missing"})
				c.Abort()
				return
			}
			tokenString = cookie
		}

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is missing"})
			c.Abort()
			return
		}

		// Parse token dengan claims
		claims := &config.JWTClaim{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Pastikan algoritma yang digunakan sesuai dengan yang diharapkan
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return config.JWT_KEY, nil
		})

		// Jika parsing token gagal
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "details": err.Error()})
			c.Abort()
			return
		}

		// Jika token tidak valid
		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is invalid"})
			c.Abort()
			return
		}

		// Simpan informasi user ke context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}
