package middleware

import (
	"github.com/DaniZGit/api.stick.it/internal/auth"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func JwtAuth() echo.MiddlewareFunc {
	// Configure middleware with the custom claims type
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(auth.JWTClaims)
		},
		SigningKey: []byte(auth.GetJwtSecret()),
		ErrorHandler: auth.JWTErrorChecker,
	}
	
	return echojwt.WithConfig(config)
}