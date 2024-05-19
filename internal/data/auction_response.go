package data

import (
	"github.com/DaniZGit/api.stick.it/internal/assetmanager"
	database "github.com/DaniZGit/api.stick.it/internal/db/generated/models"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type AuctionOffer struct {
	ID uuid.UUID `json:"id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	StartingBid  int           `json:"starting_bid"`
	Duration int `json:"duration"`
	Completed bool `json:"completed"`
	UserStickerID uuid.UUID `json:"user_sticker_id"`
	UserSticker UserSticker `json:"user_sticker"`
	LatestBid int `json:"latest_bid"`
}

type AuctionBid struct {
	ID uuid.UUID `json:"id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	Bid int `json:"bid"`
	AuctionOfferID uuid.UUID `json:"auction_offer_id"`
	UserID uuid.UUID `json:"user_id"`
	User User `json:"user"`
} 

type AuctionOfferResponse struct {
	AuctionOffer AuctionOffer `json:"auction_offer"`
}

type AuctionOffersResponse struct {
	AuctionOffers []AuctionOffer `json:"auction_offers"`
}

type AuctionBidResponse struct {
	AuctionBid AuctionBid `json:"auction_bid"`
}

type AuctionBidsResponse struct {
	AuctionBids []AuctionBid `json:"auction_bids"`
}

func CastToAuctionOfferResponse(row database.GetAuctionOfferRow) AuctionOfferResponse {
	auctionOffer := AuctionOffer{
		ID: row.ID,
		CreatedAt: row.CreatedAt,
		StartingBid: int(row.StartingBid),
		Duration: int(row.Duration),
		Completed: row.Completed,
		UserStickerID: row.UserStickerID,
		LatestBid: int(row.LatestBid),
	}

	sticker := Sticker{
		ID: uuid.NullUUID{UUID: row.StickerID, Valid: !row.StickerID.IsNil()},
		Title: row.StickerTitle,
		Type: row.StickerType,
		Top: row.StickerTop,
		Left: row.StickerLeft,
		Width: row.StickerWidth,
		Height: row.StickerHeight,
		Numerator: row.StickerNumerator,
		Denominator: row.StickerDenominator,
		Rotation: row.StickerRotation,
		PageID: row.StickerPageID,
		RarityID: row.StickerRarityID,
		CreatedAt: row.StickerCreatedAt,
		StickerID: row.StickerStickerID,
	}

	// add rarity
	if !row.StickerRarityID.UUID.IsNil() {
		sticker.Rarity = &Rarity{
			ID: row.StickerRarityID,
			Title: row.StickerRarityTitle.String,
		}
	}

	// add file
	if !row.StickerFileID.UUID.IsNil() {
		sticker.File = &File{
			ID: row.StickerFileID,
			Name: row.StickerFileName.String,
			Url: assetmanager.GetPublicAssetsFileUrl(row.StickerFilePath.String, ""),
		}
	}

	userSticker := UserSticker{
		ID: row.UserStickerID,
		UserID: row.UserStickerUserID,
		StickerID: row.UserStickerStickerID,
		Amount: int(row.UserStickerAmount),
		Sticked: row.UserStickerSticked,
		Sticker: sticker,
		Album: &Album{
			ID: row.AlbumID,
			Title: row.AlbumTitle,
		},
	}

	auctionOffer.UserSticker = userSticker
	
	return AuctionOfferResponse{
		AuctionOffer: auctionOffer,
	}
}

func CastToAuctionOffersResponse(rows []database.GetAuctionOffersRow) AuctionOffersResponse {
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
			Duration: int(row.Duration),
			Completed: row.Completed,
			UserStickerID: row.UserStickerID,
			LatestBid: int(row.LatestBid),
		}

		sticker := Sticker{
			ID: uuid.NullUUID{UUID: row.StickerID, Valid: !row.StickerID.IsNil()},
			Title: row.StickerTitle,
			Type: row.StickerType,
			Top: row.StickerTop,
			Left: row.StickerLeft,
			Width: row.StickerWidth,
			Height: row.StickerHeight,
			Numerator: row.StickerNumerator,
			Denominator: row.StickerDenominator,
			Rotation: row.StickerRotation,
			PageID: row.StickerPageID,
			RarityID: row.StickerRarityID,
			CreatedAt: row.StickerCreatedAt,
			StickerID: row.StickerStickerID,
		}

		// add rarity
		if !row.StickerRarityID.UUID.IsNil() {
			sticker.Rarity = &Rarity{
				ID: row.StickerRarityID,
				Title: row.StickerRarityTitle.String,
			}
		}

		// add file
		if !row.StickerFileID.UUID.IsNil() {
			sticker.File = &File{
				ID: row.StickerFileID,
				Name: row.StickerFileName.String,
				Url: assetmanager.GetPublicAssetsFileUrl(row.StickerFilePath.String, ""),
			}
		}

		userSticker := UserSticker{
			ID: row.UserStickerID,
			UserID: row.UserStickerUserID,
			StickerID: row.UserStickerStickerID,
			Amount: int(row.UserStickerAmount),
			Sticked: row.UserStickerSticked,
			Sticker: sticker,
			Album: &Album{
				ID: row.AlbumID,
				Title: row.AlbumTitle,
			},
		}

		auctionOffer.UserSticker = userSticker

		auctionOffers = append(auctionOffers, auctionOffer)
	}

	return AuctionOffersResponse{
		AuctionOffers: auctionOffers,
	}
}

func CastToAuctionBidResponse(row database.AuctionBid, user database.User) AuctionBidResponse {
	auctionBid := AuctionBid{
		ID: row.ID,
		CreatedAt: row.CreatedAt,
		Bid: int(row.Bid),
		AuctionOfferID: row.AuctionOfferID,
		UserID: row.UserID,
		User: User{
			ID: user.ID,
			Username: user.Username,
			Email: user.Email,
			Tokens: int(user.Tokens),
		},
	}

	return  AuctionBidResponse{
		AuctionBid: auctionBid,
	}
}

func CastToAuctionBidsResponse(rows []database.GetAuctionBidsRow) AuctionBidsResponse {
	if rows == nil || len(rows) <= 0 {
		return  AuctionBidsResponse{
			AuctionBids: []AuctionBid{},
		}
	}

	auctionBids := []AuctionBid{}
	for _, row := range rows {
		auctionBid := AuctionBid{
			ID: row.ID,
			CreatedAt: row.CreatedAt,
			Bid: int(row.Bid),
			AuctionOfferID: row.AuctionOfferID,
			UserID: row.UserID,
			User: User{
				ID: row.UserID,
				Username: row.UserUsername,
				Email: row.UserEmail,
				Tokens: int(row.UserTokens),
			},
		}

		auctionBids = append(auctionBids, auctionBid)
	}

	return  AuctionBidsResponse{
		AuctionBids: auctionBids,
	}
}