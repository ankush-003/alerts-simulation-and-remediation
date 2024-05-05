package routes

import (
	controller "Rest_server/controllers"
	"context"
	"github.com/ankush-003/alerts-simulation-and-remediation/middleware/sim/store"
	"github.com/gin-gonic/gin"
)

func PostRemedy(ctx context.Context, incomingRoutes *gin.Engine, redisClient *store.RedisStore) {
	incomingRoutes.POST("/postRemedy", controller.PostRem(ctx, redisClient))
}