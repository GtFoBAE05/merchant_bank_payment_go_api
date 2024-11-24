package test_helpers

import (
	"fmt"
	"github.com/google/uuid"
	"merchant_bank_payment_go_api/internal/entity"
	"time"
)

const CustomerFilename = "test_customers.json"
const MerchantFilename = "test_merchant.json"
const HistoryTempFilename = "test_history.json"
const PaymentTransactionTempFilename = "test_payment_transaction.json"
const BlacklistTempFilename = "test_blacklist_token.json"

var CustomerId = uuid.New()
var MerchantId = uuid.New()

var CreatedAt, _ = time.Parse("2006-01-02 15:04:05.999999999", "2024-11-22 11:31:58.769884426")
var UpdatedAt = CreatedAt

var ExpectedCustomers = []entity.Customer{{
	Id:        CustomerId,
	Username:  "budi",
	Password:  "$2a$10$2y2ss1Xs8TWZKWFS2//gnuhX/Ruhvx07lIN6jcZX1JziMvC/uLOJe",
	CreatedAt: CreatedAt,
	UpdatedAt: UpdatedAt,
}}

var ExpectedMerchants = []entity.Merchant{{
	Id:        MerchantId,
	Name:      "toko jaya",
	CreatedAt: CreatedAt,
	UpdatedAt: UpdatedAt,
}}

var ExpectedHistories = []entity.History{
	{
		Id:         uuid.New(),
		Action:     "LOGIN",
		CustomerId: CustomerId,
		Timestamp:  CreatedAt,
		Details:    "Login successful",
	},
	{
		Id:         uuid.New(),
		Action:     "PAYMENT",
		CustomerId: CustomerId,
		Timestamp:  CreatedAt.Add(1),
		Details:    fmt.Sprintf("Payment of 20000 to Merchant Id %s", MerchantId),
	},
	{
		Id:         uuid.New(),
		Action:     "LOGOUT",
		CustomerId: CustomerId,
		Timestamp:  CreatedAt.Add(2),
		Details:    "Logout successful",
	},
}

var ExpectedPayments = []entity.Payment{
	{
		Id:         uuid.New(),
		CustomerId: CustomerId,
		MerchantId: MerchantId,
		Amount:     50000,
		Timestamp:  CreatedAt,
	},
}

var ExpectedTokens = []string{"token1", "token2", "token3"}
