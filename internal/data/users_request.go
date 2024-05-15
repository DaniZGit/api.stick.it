package data

import "github.com/gofrs/uuid"

type UserRegisterParams struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UserMailConfirmationParams struct {
	Token string `json:"token" validate:"required"`
}

type UserLoginParams struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserPacksGetRequest struct {
	ID uuid.UUID `param:"id" validate:"required"`
	AlbumID uuid.UUID `query:"album_id" validate:"required"`
}

type UserStickersGetRequest struct {
	ID uuid.UUID `param:"id" validate:"required"`
	AlbumID uuid.UUID `query:"album_id"`
}

type UserPackOpenRequest struct {
	ID uuid.UUID `param:"id" validate:"required"`
	AlbumID uuid.UUID `json:"album_id" validate:"required"`
	PackID uuid.UUID `json:"pack_id" validate:"required"`
	OpenAll bool `json:"open_all"`
}

type StickUserStickerRequest struct {
	ID uuid.UUID `param:"id" validate:"required"`
	StickerID uuid.UUID `json:"sticker_id" validate:"required"`
}

type ClaimUserFreePackRequest struct {
	PackID uuid.UUID `json:"pack_id" validate:"required"`
}