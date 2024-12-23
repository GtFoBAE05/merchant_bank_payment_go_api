package route

import (
	"github.com/gin-gonic/gin"
	"merchant_bank_payment_go_api/internal/delivery/http/controller"
	"merchant_bank_payment_go_api/internal/delivery/http/middleware"
	"merchant_bank_payment_go_api/internal/usecase/impl"
)

func ConfigureRouter(router *gin.Engine, authController *controller.AuthenticationController, paymentController *controller.PaymentTransactionController, authUseCase *impl.AuthUseCaseImpl) {
	authMiddleware := middleware.AuthenticationMiddleware(authUseCase)
	publicRoute := router.Group("/api/auth")
	{
		publicRoute.POST("/login", authController.Login)
	}

	protectedRoute := router.Group("/api", authMiddleware)
	{
		protectedRoute.POST("/auth/logout", authController.Logout)
		protectedRoute.POST("/payment", paymentController.AddPayment)
	}
}
