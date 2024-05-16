package handlers

import (
	"errors"
	"net/http"

	"github.com/DaniZGit/api.stick.it/internal/app"
	"github.com/DaniZGit/api.stick.it/internal/auth"
	"github.com/DaniZGit/api.stick.it/internal/data"
	database "github.com/DaniZGit/api.stick.it/internal/db/generated/models"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
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

	// decrese user sticker amount by 1
	_, err = qtx.DecreaseUserStickerAmount(ctx.Request().Context(), a.UserStickerID)
	if err != nil {
		return ctx.ErrorResponse(http.StatusNotAcceptable, err)
	}

	// create auction offer
	auctionOffer, err := qtx.CreateAuctionOFfer(ctx.Request().Context(), database.CreateAuctionOFferParams{
		ID: uuid.Must(uuid.NewV4()),
		StartingBid: int32(a.StartingBid),
		UserStickerID: a.UserStickerID,
	})
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

////////////////////////////
/* POST - "/auction/bids" */
////////////////////////////
func CreateAuctionBid(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	a := new(data.AuctionBidCreateRequest)
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

	claims := auth.GetClaimsFromToken(*ctx)
	lastBid := 0
	lastUserID := uuid.NullUUID{}

	// get current bid
	lastAuctionBid, err := qtx.GetLatestAuctionBid(ctx.Request().Context(), a.AuctionOfferID)
	if err != nil && err != pgx.ErrNoRows {
		return ctx.ErrorResponse(http.StatusNotAcceptable, err)
	}

	// if there is no current bid, get starting_bid from auction offer instead
	if (err == pgx.ErrNoRows) {
		auctionOffer, err := qtx.GetAuctionOffer(ctx.Request().Context(), a.AuctionOfferID)
		if err != nil {
			return ctx.ErrorResponse(http.StatusNotAcceptable, err)
		}	

		lastBid = int(auctionOffer.StartingBid)
	} else {
		lastBid = int(lastAuctionBid.Bid)
		lastUserID = uuid.NullUUID{UUID: lastAuctionBid.UserID, Valid: true}
	}

	// check if current user is trying to outbid its own bid
	if !lastUserID.UUID.IsNil() && claims.UserID == lastUserID.UUID {
		return ctx.ErrorResponse(http.StatusNotAcceptable, errors.New("cannot outbid your own bid"))
	}

	// create auction bid
	auctionBid, err := qtx.CreateAuctionBid(ctx.Request().Context(), database.CreateAuctionBidParams{
		ID: uuid.Must(uuid.NewV4()),
		Bid: int32(lastBid) + 1,
		AuctionOfferID: a.AuctionOfferID,
		UserID: claims.UserID,
	})
	if err != nil {
		return ctx.ErrorResponse(http.StatusNotAcceptable, err)
	}
	
	// decrease user tokens by bid amount
	user, err := qtx.DecrementUserTokens(ctx.Request().Context(), database.DecrementUserTokensParams{
		ID: claims.UserID,
		Tokens: int64(auctionBid.Bid),
	})
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	// if there were bids before, return last bid's user its tokens
	if !lastUserID.UUID.IsNil() {
		_, err = qtx.IncrementUserTokens(ctx.Request().Context(), database.IncrementUserTokensParams{
			ID: lastUserID.UUID,
			Tokens: int64(lastBid),
		})
		if err != nil {
			return ctx.ErrorResponse(http.StatusInternalServerError, err)
		}
	}

	// commit transaction	
	err = tx.Commit(ctx.Request().Context())
	if err != nil {
		return ctx.ErrorResponse(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, data.CastToAuctionBidResponse(auctionBid, user))
}

//////////////////////////////////////
/* GET - "/auction/offers/:id/bids" */
//////////////////////////////////////
func GetAuctionBids(c echo.Context) error {
	ctx := c.(*app.ApiContext)

	a := new(data.AuctionBidsGetRequest)
	if err := ctx.Bind(a); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	if err := ctx.Validate(a); err != nil {
		return ctx.ErrorResponse(http.StatusNotImplemented, err)
	}

	auctionBids, err := ctx.Queries.GetAuctionBids(ctx.Request().Context(), a.AuctionOfferID)
	if err != nil {
		return ctx.ErrorResponse(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusCreated, data.CastToAuctionBidsResponse(auctionBids))
}
