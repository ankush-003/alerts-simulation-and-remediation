package main

import (
    "log"
    "net/http"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"

    routes "Rest_server/routes"
)

func main() {
    // Load environment variables from .env file
    err := godotenv.Load(".env")
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    // Get the port from environment variables, default to 8000 if not set
    port := os.Getenv("PORT")
    if port == "" {
        port = "8000"
    }

    // Create a new Gin router
    router := gin.New()
    router.Use(gin.Logger())

    // Enable CORS middleware
    router.Use(func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }
        c.Next()
    })

    // Define routes
    routes.AuthRoutes(router)
    routes.UserRoutes(router)
    routes.PostRemedy(router)

    // Define a simple route for testing
    router.GET("/home", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"success": "Access granted for home"})
    })

    // Run the server
    router.Run(":" + port)
}
