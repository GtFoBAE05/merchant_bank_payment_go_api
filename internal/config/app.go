package config

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"merchant_bank_payment_go_api/internal/delivery/http/controller"
	"merchant_bank_payment_go_api/internal/delivery/http/route"
	repositoryImpl "merchant_bank_payment_go_api/internal/repository/impl"
	usecaseImpl "merchant_bank_payment_go_api/internal/usecase/impl"
)

type BootstrapConfig struct {
	Log *logrus.Logger
}

func Bootstrap(config *BootstrapConfig) *gin.Engine {
	historyRepository := repositoryImpl.NewHistoryRepositoryImpl(config.Log, "internal/repository/data/History.json")
	customerRepository := repositoryImpl.NewCustomerRepository(config.Log, "internal/repository/data/Customer.json")
	merchantRepository := repositoryImpl.NewMerchantRepository(config.Log, "internal/repository/data/Merchant.json")
	authRepository := repositoryImpl.NewAuthRepository(config.Log, "internal/repository/data/BlacklistToken.json")
	paymentTransactionRepository := repositoryImpl.NewPaymentTransactionImpl(config.Log, "internal/repository/data/PaymentTransactions.json")

	historyUsecase := usecaseImpl.NewHistoryUseCaseImpl(config.Log, historyRepository)
	customerUseCase := usecaseImpl.NewCustomerUseCaseImpl(historyUsecase, customerRepository)
	merchantUseCase := usecaseImpl.NewMerchantUseCaseImpl(historyUsecase, merchantRepository)
	authUseCase := usecaseImpl.NewAuthUseCaseImpl(authRepository, customerUseCase, historyUsecase)
	paymentTransactionUseCase := usecaseImpl.NewPaymentTransactionUseCaseImpl(paymentTransactionRepository, customerUseCase,
		merchantUseCase, historyUsecase)

	authController := controller.NewAuthenticationController(config.Log, authUseCase)
	paymentController := controller.NewPaymentTransactionController(config.Log, paymentTransactionUseCase)

	router := gin.Default()
	route.ConfigureRouter(router, authController, paymentController, authUseCase)

	return router
}
