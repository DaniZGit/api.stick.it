package handlers

import (
	"net/http"

	"github.com/DaniZGit/api.stick.it/internal/app"
	"github.com/DaniZGit/api.stick.it/internal/auth"
	"github.com/DaniZGit/api.stick.it/internal/data"
	database "github.com/DaniZGit/api.stick.it/internal/db/generated/models"
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

func GetUserPacks(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	u := new(data.UserPacksGetRequest)
	if err := ctx.Bind(u); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	if err := ctx.Validate(u); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	userPacks, err := ctx.Queries.GetUserPacks(ctx.Request().Context(), database.GetUserPacksParams{
		UserID: u.ID,
		AlbumID: u.AlbumID,
	})
	if err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	return ctx.JSON(http.StatusOK, data.BuildPackResponse(userPacks, &database.File{}))
}

func GetUserStickers(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	u := new(data.UserStickersGetRequest)
	if err := ctx.Bind(u); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	if err := ctx.Validate(u); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	userStickers, err := ctx.Queries.GetUserStickers(ctx.Request().Context(), database.GetUserStickersParams{
		UserID: u.ID,
		AlbumID: u.AlbumID,
	})
	if err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	return ctx.JSON(http.StatusOK, data.BuildStickerResponse(userStickers, &database.File{}, &database.Rarity{}))
}