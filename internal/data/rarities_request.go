package data

type RarityCreateRequest struct {
	Title string `json:"title" form:"title" validate:"required"`
}