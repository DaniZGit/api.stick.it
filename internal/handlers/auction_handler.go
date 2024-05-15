package handlers

import (
	"net/http"

	"github.com/DaniZGit/api.stick.it/internal/app"
	"github.com/DaniZGit/api.stick.it/internal/data"
	database "github.com/DaniZGit/api.stick.it/internal/db/generated/models"
	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
)

//////////////////////////////
/* POST - "/auction/offers" */
//////////////////////////////
func CreateAuctionOffer(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	a := new(data.AuctionOfferCreateRequest)
	if err := ctx.Bind(a); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	if err := ctx.Validate(a); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	// start transaction
	tx, err := ctx.DBPool.Begin(ctx.Request().Context())
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}
	defer tx.Rollback(ctx.Request().Context())
	qtx := ctx.Queries.WithTx(tx)

	// create auction offer
	auctionOffer, err := qtx.CreateAuctionOFfer(ctx.Request().Context(), database.CreateAuctionOFferParams{
		Column1: uuid.Must(uuid.NewV4()),
		Column2: int32(a.StartingBid),
		Column3: a.UserStickerID,
	})
	if err != nil {
		return ctx.ErrorResponse(http.StatusNotAcceptable, err)
	}

	// decrese user sticker amount by 1
	_, err = qtx.DecreaseUserStickerAmount(ctx.Request().Context(), a.UserStickerID)
	if err != nil {
		return ctx.ErrorResponse(http.StatusNotAcceptable, err)
	}

	// commit transaction	
	err = tx.Commit(ctx.Request().Context())
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, data.CastToAuctionOfferResponse(auctionOffer))
}

/////////////////////////////
/* GET - "/auction/offers" */
/////////////////////////////
func GetAuctionOffers(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	auctionOffers, err := ctx.Queries.GetAuctionOffers(ctx.Request().Context())
	if err != nil {
		return ctx.ErrorResponse(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusCreated, data.CastToAuctionOffersResponse(auctionOffers))
}