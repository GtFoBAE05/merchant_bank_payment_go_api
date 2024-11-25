package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"merchant_bank_payment_go_api/internal/model"
	"merchant_bank_payment_go_api/internal/usecase"
	auth "merchant_bank_payment_go_api/internal/utils"
	"net/http"
	"strings"
)

func AuthenticationMiddleware(authUseCase usecase.AuthUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		logrus.Infof("Starting Authorization header validation for %s", c.Request.URL.Path)

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logrus.Warn("Missing Authorization header")
			c.JSON(http.StatusUnauthorized, model.CommonResponse[interface{}]{
				HttpStatus: http.StatusUnauthorized,
				Message:    "Authorization Header is required",
				Data:       nil,
			})
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			logrus.Warn("Invalid Authorization header format")
			c.JSON(http.StatusUnauthorized, model.CommonResponse[interface{}]{
				HttpStatus: http.StatusUnauthorized,
				Message:    "Invalid Authorization Header format, must be 'Bearer <token>'",
				Data:       nil,
			})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]
		logrus.Debugf("Verifying token for request: %s", c.Request.URL.Path)

		valid, err := auth.VerifyAccessToken(tokenString)
		if err != nil || !valid {
			if err != nil {
				logrus.Errorf("Error verifying token: %v", err)
			} else {
				logrus.Warn("Expired or invalid token detected")
			}

			c.JSON(http.StatusUnauthorized, model.CommonResponse[interface{}]{
				HttpStatus: http.StatusUnauthorized,
				Message:    "Invalid or expired token",
				Data:       nil,
			})
			c.Abort()
			return
		}

		isBlacklisted, err := authUseCase.IsTokenBlacklisted(tokenString)
		if err != nil {
			logrus.Errorf("Error checking blacklist status: %v", err)
			c.JSON(http.StatusInternalServerError, model.CommonResponse[interface{}]{
				HttpStatus: http.StatusInternalServerError,
				Message:    "Internal server error",
				Data:       nil,
			})
			c.Abort()
			return
		}

		if isBlacklisted {
			logrus.Warnf("Token %s is already blacklisted", tokenString)
			c.JSON(http.StatusForbidden, model.CommonResponse[interface{}]{
				HttpStatus: http.StatusForbidden,
				Message:    "Token is already blacklisted",
				Data:       nil,
			})
			c.Abort()
			return
		}

		userId, err := auth.ExtractIDFromToken(tokenString)
		if err != nil {
			logrus.Errorf("Error extracting ID from token: %v", err)
			c.JSON(http.StatusUnauthorized, model.CommonResponse[interface{}]{
				HttpStatus: http.StatusUnauthorized,
				Message:    "Invalid token data",
				Data:       nil,
			})
			c.Abort()
			return
		}

		c.Set("user_id", userId)
		c.Set("token", tokenString)
		c.Next()
	}
}
