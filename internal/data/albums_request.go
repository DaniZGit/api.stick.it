package data

import "github.com/gofrs/uuid"

type AlbumCreateRequest struct {
	Title string `json:"title" form:"title" validate:"required"`
	DateFrom string `json:"date_from" form:"date_from" validate:"required"`
	DateTo string `json:"date_to" form:"date_to" validate:"required"`
	Featured bool `json:"featured" form:"featured"`
}

type AlbumUpdateRequest struct {
	ID uuid.UUID `param:"id" validate:"required"`
	Title string `json:"title" form:"title"`
	DateFrom string `json:"date_from" form:"date_from"`
	DateTo string `json:"date_to" form:"date_to"`
	Featured bool `json:"featured" form:"featured"`
	FileID string `json:"file_id" form:"file_id"`
}

type AlbumDeleteRequest struct {
	ID string `param:"id" validate:"required"`
}