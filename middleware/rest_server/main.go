package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/ankush-003/alerts-simulation-and-remediation/middleware/sim/store"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	routes "Rest_server/routes"
)

func main() {
	log.Println("Starting the server")

	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}

	logger := log.New(os.Stdout, "Rest Server:", log.LstdFlags)

	// Get the port from environment variables, default to 8000 if not set
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
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

	// Create a new Gin router
	router := gin.New()
	router.Use(gin.Logger())

	// Enable CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS ,HEAD, PUT, FETCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Authorization, Access-Control-Request-Method, Access-Control-Request-Headers")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// Define routes
	routes.PostRemedy(ctx, router, redis)
	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	// Define a simple route for testing
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": "Welcome to the alerts simulation and remediation server!"})
	})

	// Run the server
	router.Run(":" + port)
}
