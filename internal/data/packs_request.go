package data

import (
	"github.com/gofrs/uuid"
)

type PackCreateRequest struct {
	AlbumID uuid.UUID `json:"album_id" form:"album_id" validate:"required"`
	Title string `json:"title" form:"title" validaiton:"required"`
	Price int `json:"price" form:"price" validaiton:"required"`
	Amount int `json:"amount" form:"amount" validaiton:"required"`
}

type PackUpdateRequest struct {
	ID uuid.UUID `param:"id" validate:"required"`
	Title string `json:"title" form:"title" validaiton:"required"`
	Price int `json:"price" form:"price" validaiton:"required"`
	Amount int `json:"amount" form:"amount" validaiton:"required"`
	FileID string `json:"file_id" form:"file_id"`
}

type PackDeleteRequest struct {
	ID uuid.UUID `param:"id" validate:"required"`
}

type AlbumPacksGetRequest struct {
	AlbumID uuid.UUID `param:"id" validate:"required"`
}

type PackRaritiesGetRequest struct {
	ID uuid.UUID `param:"id" validate:"required"`
}

type PackRarityCreateRequest struct {
	PackID uuid.UUID `json:"pack_id" form:"pack_id" validate:"required"`
	RarityID uuid.UUID `json:"rarity_id" form:"rarity_id" validate:"required"`
	DropChance float32 `json:"drop_chance" form:"drop_chance" validate:"required"`
}

type PackRarityUpdateRequest struct {
	ID uuid.UUID `param:"id" validate:"required"`
	DropChance float32 `json:"drop_chance" form:"drop_chance" validate:"required"`
}

type PackRarityDeleteRequest struct {
	ID uuid.UUID `param:"id" validate:"required"`
}