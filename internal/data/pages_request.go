package data

import "github.com/gofrs/uuid"

type PageCreateRequest struct {
	AlbumID uuid.UUID `json:"album_id" form:"album_id" validate:"required"`
	SortOrder int32 `json:"sort_order" form:"sort_order"`
}

type PageGetRequest struct {
	ID string `param:"id" validate:"required"`
}

type PageUpdateRequest struct {
	ID uuid.UUID `param:"id" validate:"required"`
	SortOrder int `json:"sort_order" form:"sort_order"`
	FileID string `json:"file_id" form:"file_id"`
}

type PageDeleteRequest struct {
	ID uuid.UUID `param:"id" validate:"required"`
}