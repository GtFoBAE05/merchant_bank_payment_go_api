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
	authRepository := repositoryImpl.NewAuthRepository(config.Log, "internal/repository/data/BlacklistToken.json")

	historyUsecase := usecaseImpl.NewHistoryUseCaseImpl(config.Log, historyRepository)
	customerUseCase := usecaseImpl.NewCustomerUseCaseImpl(config.Log, customerRepository)
	authUseCase := usecaseImpl.NewAuthUseCaseImpl(config.Log, authRepository, customerUseCase, historyUsecase)

	authController := controller.NewAuthController(config.Log, authUseCase)

	router := gin.Default()
	route.ConfigureRouter(router, authController, authUseCase)

	return router
}
