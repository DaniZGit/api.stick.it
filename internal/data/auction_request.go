package data

import "github.com/gofrs/uuid"

type AuctionOfferCreateRequest struct {
	UserStickerID uuid.UUID `json:"user_sticker_id" validate:"required"`
	StartingBid int `json:"starting_bid" validate:"required"`
}

type AuctionBidCreateRequest struct {
	AuctionOfferID uuid.UUID `json:"auction_offer_id" validate:"required"`
}

type AuctionBidsGetRequest struct {
	AuctionOfferID uuid.UUID `query:"auction_offer_id" validate:"required"`
}