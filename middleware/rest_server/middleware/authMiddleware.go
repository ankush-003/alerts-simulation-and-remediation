package middleware

import (
	helper "Rest_server/helpers"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("Authorization")
		// fmt.Println(c.Request.Header)
		// fmt.Println(c.Request)
		// fmt.Println("Here", clientToken)
		// fmt.Println(clientToken)
		// clientToken = strings.Split(clientToken, "Bearer ")[1]
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("No Authorization header provided")})
			c.Abort()
			return
		}
		clientToken = clientToken[7:]

		fmt.Println("He")

		claims, err := helper.ValidateToken(clientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		fmt.Println(claims)

		c.Set("email", claims.Email)
		c.Set("first_name", claims.First_name)
		c.Set("last_name", claims.Last_name)
		c.Set("uid", claims.Uid)
		c.Set("user_type", claims.User_type)
		c.Next()

		// fmt.Println("end")

	}
}

func GetUserIdContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Attempt to retrieve the user ID from the context set by the authentication middleware
		userID, exists := c.Get("uid")
		if !exists {
			// If user ID is not found in the context, attempt to extract it from the request headers
			userIDFromHeader := c.GetHeader("uid")
			if userIDFromHeader == "" {
				// If user ID is not found in the headers, return an empty string or handle as per your application logic
				c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found"})
				c.Abort()
				return
			}
			c.Set("uid", userIDFromHeader) // Set the user ID in the context for future use
			return
		}

		// Convert the user ID to a string
		userIDString, ok := userID.(string)
		if !ok {
			// Handle error if user ID is not in the expected format
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
			c.Abort()
			return
		}

		// Set the user ID in the context for future use
		c.Set("uid", userIDString)

		c.Next() // Proceed to the next middleware or route handler
	}
}
