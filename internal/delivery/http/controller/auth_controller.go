package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"merchant_bank_payment_go_api/internal/usecase"
	"net/http"
)

type AuthController struct {
	Log         *logrus.Logger
	AuthUseCase usecase.AuthUseCaseInterface
}

func NewAuthController(logger *logrus.Logger, authUseCase usecase.AuthUseCaseInterface) *AuthController {
	return &AuthController{
		Log:         logger,
		AuthUseCase: authUseCase,
	}
}

func (ac *AuthController) Login(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}

func (ac *AuthController) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}
