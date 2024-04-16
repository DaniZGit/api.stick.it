package auth

import (
	"net/http"
	"time"

	"github.com/DaniZGit/api.stick.it/environment"
	"github.com/DaniZGit/api.stick.it/internal/app"
	database "github.com/DaniZGit/api.stick.it/internal/db/generated/models"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type JWTClaims struct {
	UserID  uuid.UUID `json:"user_id"`
	RoleID uuid.NullUUID  `json:"role_id"`
	jwt.RegisteredClaims
}

func GetJwtSecret() string {
	return environment.JwtSecret()
}

func GetJwtTokenName() string {
	return "access-token"
}

/* Generates a new jwt token */
func CreateJwtToken(user database.User) (string, error) {
	// set expiration time for token and cookie
	expirationTime := time.Now().Add(72 * time.Hour)

	// generate jwt token
	claims := &JWTClaims{
		UserID: user.ID,
		RoleID: user.RoleID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID: user.ID.String(),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// generate encoded token.
	t, err := token.SignedString([]byte(GetJwtSecret()))
	if err != nil {
		return "", err
	}

	return t, nil
}

func GetClaimsFromToken(ctx app.ApiContext) *JWTClaims {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(*JWTClaims)

	return claims;
}

// JWTErrorChecker will be executed when user try to access a protected path.
func JWTErrorChecker(ctx echo.Context, err error) error {
	return ctx.JSON(http.StatusUnauthorized, echo.Map{
		"error": "Unaothurazioned",
	})
}