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
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/paymentintent"
)

//////////////////////////////////////////////////
/* POST - "/transactions/create-payment-intent" */
//////////////////////////////////////////////////
func CreatePaymentIntent(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	t := new(data.TransactionCreatePaymentIntentRequest)
	if err := ctx.Bind(t); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	if err := ctx.Validate(t); err != nil {
		return ctx.ErrorResponse(http.StatusUnprocessableEntity, err)
	}

	bundle, err := ctx.Queries.GetBundle(ctx.Request().Context(), t.BundleID)
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	price, err := bundle.Price.Float64Value()
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	params := stripe.PaymentIntentParams{
		Amount: stripe.Int64(int64(price.Float64 * 100)),
		Currency: t.Currency,
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}
	paymentIntent, err := paymentintent.New(&params)
	if err != nil {
		return ctx.ErrorResponse(http.StatusUnprocessableEntity, err)
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"client_secret": paymentIntent.ClientSecret,
	})
}

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

	t := new(data.TransactionBundleBuyRequest)
	if err := ctx.Bind(t); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	if err := ctx.Validate(t); err != nil {
		return ctx.ErrorResponse(http.StatusUnprocessableEntity, err)
	}

	bundle, err := ctx.Queries.GetBundle(ctx.Request().Context(), t.BundleID)
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	claims := auth.GetClaimsFromToken(*ctx)
	user, err := ctx.Queries.IncrementUserTokens(ctx.Request().Context(), database.IncrementUserTokensParams{
		ID: claims.UserID,
		Tokens: int64(bundle.Tokens + bundle.Bonus),
	})
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"tokens": user.Tokens,
	})
}