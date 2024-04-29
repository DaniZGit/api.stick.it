package data

import "github.com/gofrs/uuid"

type TransactionCreatePaymentIntentRequest struct {
	BundleID uuid.UUID `json:"bundle_id" form:"bundle_id" validate:"required"`
	Currency *string `json:"currency" validate:"required"`
}

type TransactionPackBuyRequest struct {
	PackID uuid.UUID `json:"pack_id" form:"pack_id" validate:"required"`
	Amount int `json:"amount" validate:"required"`
}

type TransactionBundleBuyRequest struct {
	BundleID uuid.UUID `json:"bundle_id" form:"bundle_id" validate:"required"`
}