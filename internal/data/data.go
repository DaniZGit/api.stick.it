package data

type Metadata struct {
	CurrPage int32 `json:"curr_page"`
	PageSize int32 `json:"page_size"`
	TotalRecords int32 `json:"total_records"`
	FirstPage int32 `json:"first_page"`
	LastPage int32 `json:"last_page"`
}