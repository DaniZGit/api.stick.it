package data

import "github.com/gofrs/uuid"

type AlbumCreateRequest struct {
	Title string `json:"title" form:"title" validate:"required"`
	DateFrom string `json:"date_from" form:"date_from" validate:"required"`
	DateTo string `json:"date_to" form:"date_to" validate:"required"`
	Featured bool `json:"featured" form:"featured"`
	PageNumerator int `json:"page_numerator" form:"page_numerator" validate:"required"` 
	PageDenominator int `json:"page_denominator" form:"page_denominator" validate:"required"` 
}

type AlbumGetRequest struct {
	ID string `param:"id" validate:"required"`
}

type AlbumUpdateRequest struct {
	ID uuid.UUID `param:"id" validate:"required"`
	Title string `json:"title" form:"title"`
	DateFrom string `json:"date_from" form:"date_from"`
	DateTo string `json:"date_to" form:"date_to"`
	Featured bool `json:"featured" form:"featured"`
	PageNumerator int `json:"page_numerator" form:"page_numerator" validate:"required"` 
	PageDenominator int `json:"page_denominator" form:"page_denominator" validate:"required"` 
	FileID string `json:"file_id" form:"file_id"`
}

type AlbumDeleteRequest struct {
	ID string `param:"id" validate:"required"`
}

type AlbumPacksGetRequest struct {
	AlbumID uuid.UUID `param:"id" validate:"required"`
}

type AlbumPagesGetRequest struct {
	AlbumID uuid.UUID `param:"id" validate:"required"`
	From int `query:"from" validate:"required"`
	To int `query:"to" validate:"required"`
}
