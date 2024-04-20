package routes

import (
	controller "Rest_server/controllers"
	"Rest_server/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine){

	incomingRoutes.POST("users/signup", controller.Signup())


	incomingRoutes.POST("users/login", controller.Login())
	authenticated := incomingRoutes.Group("/")
    authenticated.Use(middleware.Authenticate())
    {
        authenticated.POST("users/alertconfig", controller.AlertConfig())
    }

}
