package handlers

import (
	"errors"
	"net/http"

	"github.com/DaniZGit/api.stick.it/internal/app"
	"github.com/DaniZGit/api.stick.it/internal/auth"
	"github.com/DaniZGit/api.stick.it/internal/data"
	database "github.com/DaniZGit/api.stick.it/internal/db/generated/models"
	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
)

/////////////////////////////////////
/* POST - "/transactions/pack/:id" */
/////////////////////////////////////
func BuyPack(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	t := new(data.TransactionPackBuyRequest)
	if err := ctx.Bind(t); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	if err := ctx.Validate(t); err != nil {
		return ctx.ErrorResponse(http.StatusUnprocessableEntity, err)
	}

	claims := auth.GetClaimsFromToken(*ctx)
	user, err := ctx.Queries.GetUserByID(ctx.Request().Context(), claims.UserID)
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	pack, err := ctx.Queries.GetPack(ctx.Request().Context(), t.PackID)
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	// check if user has enough tokens to buy the pack
	if (user.Tokens < int64(pack.Price) * int64(t.Amount)) {
		return ctx.ErrorResponse(http.StatusUnprocessableEntity, errors.New("user doesn't have enough tokens"))
	}

	// decrement user tokens by pack cost * amount of packs bought
	user, err = ctx.Queries.DecrementUserTokens(ctx.Request().Context(), database.DecrementUserTokensParams{
		ID: user.ID,
		Tokens: int64(pack.Price) * int64(t.Amount),
	})
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	// buy pack
	boughtPack, err := ctx.Queries.CreateUserPack(ctx.Request().Context(), database.CreateUserPackParams{
		ID: uuid.Must(uuid.NewV4()),
		UserID: user.ID,
		PackID: pack.ID,
		Amount: int32(t.Amount),
	})
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, data.BuildPackResponse(boughtPack, &database.File{}))
}

func BuyBundle(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	claims := auth.GetClaimsFromToken(*ctx)
	user, err := ctx.Queries.IncrementUserTokens(ctx.Request().Context(), database.IncrementUserTokensParams{
		ID: claims.UserID,
		Tokens: 1000,
	})
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"user": user,
	})
}