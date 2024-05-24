package routes

import (
	controller "Rest_server/controllers"
	"Rest_server/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {

	incomingRoutes.POST("users/signup", controller.Signup())

	incomingRoutes.POST("users/login", controller.Login())
	authenticated := incomingRoutes.Group("/")
	// authenticated.GET("/alerts", controller.GetAllAlerts())

	authenticated.Use(middleware.Authenticate())
	{
		authenticated.GET("/alerts", controller.GetUserAlerts())
		authenticated.PUT("/acknowledge", controller.AcknowledgeLog)
		authenticated.POST("users/alertconfig", controller.AlertConfig())
	}

}
