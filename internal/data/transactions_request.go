package data

import "github.com/gofrs/uuid"

type TransactionPackBuyRequest struct {
	PackID uuid.UUID `json:"pack_id" form:"pack_id" validate:"required"`
	Amount int `json:"amount" validate:"required"`
}
