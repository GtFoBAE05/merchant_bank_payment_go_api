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
	customerRepository := repositoryImpl.NewCustomerRepository(config.Log, "internal/repository/data/Customer.json")
	authRepository := repositoryImpl.NewAuthRepository(config.Log, "internal/repository/data/BlacklistToken.json")

	customerUseCase := usecaseImpl.NewCustomerUseCase(config.Log, customerRepository)
	authUseCase := usecaseImpl.NewAuthUseCase(config.Log, authRepository, customerUseCase)

	authController := controller.NewAuthController(config.Log, authUseCase)

	router := gin.Default()
	route.ConfigureRouter(router, authController, authUseCase)

	return router
}
