package config

import "github.com/golang-jwt/jwt/v5"

var JWT_KEY = []byte("d09ed564035ad789557879972cb2f78e3595b19fad6edf0924c26b8571f11311")

type JWTClaim struct {
	Username string
	jwt.RegisteredClaims
}
