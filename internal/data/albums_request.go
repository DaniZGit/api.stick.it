package data

type AlbumCreateRequest struct {
	Title string `json:"title" form:"title" validate:"required"`
	DateFrom string `json:"date_from" form:"date_from" validate:"required"`
	DateTo string `json:"date_to" form:"date_to" validate:"required"`
}

type AlbumDeleteRequest struct {
	Title string `param:"title" validate:"required"`
}