package handlers

import (
	"net/http"

	"github.com/DaniZGit/api.stick.it/internal/app"
	"github.com/DaniZGit/api.stick.it/internal/auth"
	"github.com/labstack/echo/v4"
)

////////////////////////
/* GET - "/users/:id" */
////////////////////////
func GetUser(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	claims := auth.GetClaimsFromToken(*ctx)

	return ctx.JSON(http.StatusOK, echo.Map{
		"user_id": claims.UserID,
		"role_id": claims.RoleID,
	})
}