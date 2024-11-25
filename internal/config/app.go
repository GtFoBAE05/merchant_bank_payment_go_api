package config

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"merchant_bank_payment_go_api/internal/delivery/http/controller"
	"merchant_bank_payment_go_api/internal/delivery/http/route"
	repositoryImpl "merchant_bank_payment_go_api/internal/repository/impl"
	usecaseImpl "merchant_bank_payment_go_api/internal/usecase/impl"
)

func Bootstrap(logger *logrus.Logger) *gin.Engine {
	historyRepository := repositoryImpl.NewHistoryRepositoryImpl(logger, "internal/repository/data/History.json")
	customerRepository := repositoryImpl.NewCustomerRepositoryImpl(logger, "internal/repository/data/Customer.json")
	merchantRepository := repositoryImpl.NewMerchantRepositoryImpl(logger, "internal/repository/data/Merchant.json")
	authRepository := repositoryImpl.NewAuthRepository(logger, "internal/repository/data/BlacklistToken.json")
	paymentTransactionRepository := repositoryImpl.NewPaymentTransactionImpl(logger, "internal/repository/data/PaymentTransactions.json")

	historyUsecase := usecaseImpl.NewHistoryUseCaseImpl(logger, historyRepository)
	customerUseCase := usecaseImpl.NewCustomerUseCaseImpl(historyUsecase, customerRepository)
	merchantUseCase := usecaseImpl.NewMerchantUseCaseImpl(historyUsecase, merchantRepository)
	authUseCase := usecaseImpl.NewAuthUseCaseImpl(authRepository, customerUseCase, historyUsecase)
	paymentTransactionUseCase := usecaseImpl.NewPaymentTransactionUseCaseImpl(paymentTransactionRepository, customerUseCase,
		merchantUseCase, historyUsecase)

	authController := controller.NewAuthenticationController(logger, authUseCase)
	paymentController := controller.NewPaymentTransactionController(logger, paymentTransactionUseCase)

	router := gin.Default()
	route.ConfigureRouter(router, authController, paymentController, authUseCase)

	return router
}
