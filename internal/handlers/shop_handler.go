package handlers

import (
	"net/http"

	"github.com/DaniZGit/api.stick.it/internal/app"
	"github.com/DaniZGit/api.stick.it/internal/data"
	"github.com/labstack/echo/v4"
)

/////////////////////////
/* GET - "/shop/packs" */
/////////////////////////
func GetShopPacks(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	packs, err := ctx.Queries.GetShopPacks(ctx.Request().Context())
	if err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	return ctx.JSON(http.StatusOK, data.BuildShopResponse(packs))
}