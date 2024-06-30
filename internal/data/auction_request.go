package data

import "github.com/gofrs/uuid"

type AuctionOfferCreateRequest struct {
	UserStickerID uuid.UUID `json:"user_sticker_id" validate:"required"`
	StartingBid int `json:"starting_bid" validate:"required"`
}

type AuctionOffersGetRequest struct {
	Limit int `query:"limit" validate:"required"`
	Page *int `query:"page" validate:"required"`
	SortField string `query:"sort_field" validate:"required,oneof=bid timespan"`
	SortOrder string `query:"sort_order" validate:"required,oneof=asc ASC desc DESC"`
	AlbumID string `query:"album_id"`
}

type AuctionBidCreateRequest struct {
	AuctionOfferID uuid.UUID `json:"auction_offer_id" validate:"required"`
	Bid int `json:"bid" validate:"required"`
}

type AuctionBidsGetRequest struct {
	AuctionOfferID uuid.UUID `query:"auction_offer_id" validate:"required"`
}