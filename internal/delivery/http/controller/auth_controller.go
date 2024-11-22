package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"merchant_bank_payment_go_api/internal/model"
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
	var loginRequest model.LoginRequest
	ac.Log.Debug("Attempting login for user")

	err := c.ShouldBind(&loginRequest)
	if err != nil {
		ac.Log.Errorf("Invalid login request: %v", err)
		c.JSON(http.StatusBadRequest, model.CommonResponse[interface{}]{
			HttpStatus: http.StatusBadRequest,
			Message:    "Invalid body request",
			Data:       nil,
		})
		return
	}

	token, err := ac.AuthUseCase.Login(loginRequest)
	if err != nil {
		ac.Log.Errorf("Login failed for user %s: %v", loginRequest.Username, err)
		c.JSON(http.StatusUnauthorized, model.CommonResponse[interface{}]{
			HttpStatus: http.StatusUnauthorized,
			Message:    err.Error(),
			Data:       nil,
		})
		return
	}

	ac.Log.Infof("Successful login for user: %s", loginRequest.Username)
	c.JSON(http.StatusOK, model.CommonResponse[model.LoginResponse]{
		HttpStatus: http.StatusOK,
		Message:    "Successfully logged in",
		Data:       token,
	})
}

func (ac *AuthController) Logout(c *gin.Context) {
	token, exists := c.Get("token")
	if !exists {
		ac.Log.Warn("Token not found in context")
		c.JSON(http.StatusUnauthorized, model.CommonResponse[interface{}]{
			HttpStatus: http.StatusUnauthorized,
			Message:    "Token not found",
			Data:       nil,
		})
		return
	}

	tokenString, _ := token.(string)
	err := ac.AuthUseCase.Logout(tokenString)
	if err != nil {
		ac.Log.Warn("Error when logout")
		c.JSON(http.StatusUnauthorized, model.CommonResponse[interface{}]{
			HttpStatus: http.StatusUnauthorized,
			Message:    err.Error(),
			Data:       nil,
		})
		return
	}

	ac.Log.Infof("Successfully logged out and blacklisted token for user: %s", tokenString)
	c.JSON(http.StatusOK, model.CommonResponse[interface{}]{
		HttpStatus: http.StatusOK,
		Message:    "Successfully logged out",
		Data:       nil,
	})
}
