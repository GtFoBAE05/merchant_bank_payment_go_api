package model

type PaymentRequest struct {
	MerchantId string `json:"merchantId" binding:"required"`
	Amount     int64  `json:"amount" binding:"required"`
}
