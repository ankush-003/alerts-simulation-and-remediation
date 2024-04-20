package main

import(
	routes "Rest_server/routes"
	"os"
	"log"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"context"
	"github.com/ankush-003/alerts-simulation-and-remediation/middleware/sim/store"
)

func main(){
	err := godotenv.Load(".env")

	logger := log.New(os.Stdout, "Rest Server:", log.LstdFlags)

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")

	if port==""{
		port="8000"
	}

	redis_addr := os.Getenv("REDIS_ADDR")
	if redis_addr == "" {
		logger.Println("REDIS_ADDR not set, using default localhost:6379")
		redis_addr = "localhost:6379"
	}

	ctx := context.Background()

	redis, redisErr := store.NewRedisStore(ctx, redis_addr)

	if redisErr != nil {
		logger.Fatalf("Error creating redis store: %s\n", redisErr)
	}

	defer redis.Close()

	router := gin.New()
	router.Use(gin.Logger())

	//router.Static("/", "./dashboard")

	// routes.AuthRoutes(router)
	// routes.UserRoutes(router)
	routes.PostRemedy(ctx, router, redis)
	
	//routes.AlertConfigRoutes(router)

	router.GET("/home", func(c *gin.Context){
		c.JSON(200, gin.H{"success":"Access granted for home"})
	})

	/*router.GET("/alertConfig", func(c *gin.Context){
		c.JSON(200, gin.H{"success":"Access granted for alert-config"})
	})*/

	router.Run(":" + port)
}	
