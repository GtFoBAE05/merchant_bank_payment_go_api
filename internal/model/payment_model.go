package model

type PaymentRequest struct {
	MerchantId string `json:"merchantId" binding:"required" validation:"min=1"`
	Amount     int64  `json:"amount" binding:"required"`
}
