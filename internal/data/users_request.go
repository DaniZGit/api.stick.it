package data

import "github.com/gofrs/uuid"

type UserRegisterParams struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
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
	AlbumID uuid.UUID `query:"album_id" validate:"required"`
}