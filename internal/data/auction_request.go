package data

import "github.com/gofrs/uuid"

type AuctionOfferCreateRequest struct {
	UserStickerID uuid.UUID `json:"user_sticker_id" validate:"required"`
	StartingBid int `json:"starting_bid" validate:"required"`
}