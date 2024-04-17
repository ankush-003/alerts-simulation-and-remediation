package main

import(
	routes "Rest_server/routes"
	"os"
	"log"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main(){
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")

	if port==""{
		port="8000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	//router.Static("/", "./dashboard")

	// routes.AuthRoutes(router)
	// routes.UserRoutes(router)
	routes.PostRemedy(router)
	
	//routes.AlertConfigRoutes(router)

	router.GET("/home", func(c *gin.Context){
		c.JSON(200, gin.H{"success":"Access granted for home"})
	})

	/*router.GET("/alertConfig", func(c *gin.Context){
		c.JSON(200, gin.H{"success":"Access granted for alert-config"})
	})*/

	router.Run(":" + port)
}	
