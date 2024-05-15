package data

import (
	database "github.com/DaniZGit/api.stick.it/internal/db/generated/models"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type AuctionOffer struct {
	ID uuid.UUID `json:"id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	StartingBid  int           `json:"starting_bid"`
	UserStickerID uuid.UUID `json:"user_sticker_id"`
}

type AuctionOfferResponse struct {
	AuctionOffer AuctionOffer `json:"auction_offer"`
}

type AuctionOffersResponse struct {
	AuctionOffers []AuctionOffer `json:"auction_offers"`
}

func CastToAuctionOfferResponse(row database.AuctionOffer) AuctionOfferResponse {
	return AuctionOfferResponse{
		AuctionOffer: AuctionOffer{
			ID: row.ID,
			CreatedAt: row.CreatedAt,
			StartingBid: int(row.StartingBid),
			UserStickerID: row.UserStickerID,
		},
	}
}

func CastToAuctionOffersResponse(rows []database.AuctionOffer) AuctionOffersResponse {
	if rows == nil || len(rows) <= 0 {
		return  AuctionOffersResponse{
			AuctionOffers: []AuctionOffer{},
		}
	}

	auctionOffers := []AuctionOffer{}
	for _, row := range rows {
		auctionOffer := AuctionOffer{
			ID: row.ID,
			CreatedAt: row.CreatedAt,
			StartingBid: int(row.StartingBid),
			UserStickerID: row.UserStickerID,
		}

		auctionOffers = append(auctionOffers, auctionOffer)
	}

	return AuctionOffersResponse{
		AuctionOffers: auctionOffers,
	}
}