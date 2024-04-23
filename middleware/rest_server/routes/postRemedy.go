package routes

import (
	controller "Rest_server/controllers"

	"github.com/gin-gonic/gin"
)

func PostRemedy(incomingRoutes *gin.Engine){
	incomingRoutes.POST("/postRemedy", controller.PostRem())
}