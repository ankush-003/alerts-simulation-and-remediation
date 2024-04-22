package controllers

import (
	"Rest_server/database"
	helper "Rest_server/helpers"
	middleware "Rest_server/middleware"
	"Rest_server/models"
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection

func init() {
	userCollection = database.OpenCollection("AlertSimAndRemediation","Users")
}
var alertCollection *mongo.Collection

func init() {
	alertCollection = database.OpenCollection("AlertSimAndRemediation","Users")
}


var validate = validator.New()

func HashPassword(password string) string{
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err!=nil{
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string)(bool, string){
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err!= nil {
		msg = fmt.Sprintf("email of password is incorrect")
		check=false
	}
	return check, msg
}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a context with a timeout for database operations
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel() // Always defer the cancel function to release resources

		// Bind the request JSON to a User struct
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Perform validation of the user struct
		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// Check if the email already exists in the database
		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while checking for the email"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this email already exists"})
			return
		}

		// Check if the phone number already exists in the database
		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while checking for the phone number"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this phone number already exists"})
			return
		}

		// Hash the password before storing it in the database
		password := HashPassword(*user.Password)
		user.Password = &password

		// Set the created_at and updated_at fields
		user.Created_at = time.Now()
		user.Updated_at = time.Now()

		// Generate authentication tokens
		token, refreshToken, _ := helper.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, user.ID.Hex())
		user.Token = &token
		user.Refresh_token = &refreshToken

		// Insert only the required fields into the database
		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, bson.M{
    		"first_name":    user.First_name,
    		"last_name":     user.Last_name,
    		"email":         user.Email,
    		"password":      user.Password,
    		"phone":         user.Phone,
    		"user_type":     user.User_type,
    		"created_at":    user.Created_at,
    		"updated_at":    user.Updated_at,
	})
	if insertErr != nil {
    	msg := fmt.Sprintf("User item was not created")
    	c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
    	return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour), // Session expires in 24 hours
		HttpOnly: true,
	})

	
	// Return success response with the inserted document ID
	//c.JSON(http.StatusOK, gin.H{"inserted_id": resultInsertionNumber.InsertedID})
	c.JSON(http.StatusOK, gin.H{"token": token, "inserted_id": resultInsertionNumber.InsertedID})
	}
}


func Login() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return 
		}



		err := userCollection.FindOne(ctx, bson.M{"email":user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"email or password is incorrect"})
			return
		}


		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if passwordIsValid != true{
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if foundUser.Email == nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"user not found"})
		}
		token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, *foundUser.User_type, foundUser.User_id)
		helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)
		err = userCollection.FindOne(ctx, bson.M{"user_id":foundUser.User_id}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		/*type Session struct {
			Token    string
			ExpireAt time.Time
		}

		var sessions map[string]Session

		session := Session{
			Token:    token,
			ExpireAt: time.Now().Add(24 * time.Hour), // Session expires in 24 hours
		}

		// Store session data
		sessions[token] = session

		// Set session cookie in the response
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "session_token",
			Value:    token,
			Expires:  session.ExpireAt,
			HttpOnly: true,
		})*/

		/*http.SetCookie(c.Writer, &http.Cookie{
			Name:     "session_token",
			Value:    token,
			Expires:  time.Now().Add(24 * time.Hour), // Session expires in 24 hours
			HttpOnly: true,
		})*/

		//c.JSON(http.StatusOK, foundUser)
		// Store the tokens in the session storage
		//sessionStorage.setItem('accessToken', token);
		c.JSON(http.StatusOK, gin.H{"token": token, "user": foundUser})
	}
}

func GetUsers() gin.HandlerFunc{
	return func(c *gin.Context){
		if err := helper.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		
		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage <1{
			recordPerPage = 10
		}
		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 !=nil || page<1{
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{"$match", bson.D{{}}}}
		groupStage := bson.D{{"$group", bson.D{
			{"_id", bson.D{{"_id", "null"}}}, 
			{"total_count", bson.D{{"$sum", 1}}}, 
			{"data", bson.D{{"$push", "$$ROOT"}}}}}}
		projectStage := bson.D{
			{"$project", bson.D{
				{"_id", 0},
				{"total_count", 1},
				{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},}}}
result,err := userCollection.Aggregate(ctx, mongo.Pipeline{
	matchStage, groupStage, projectStage})
defer cancel()
if err!=nil{
	c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while listing user items"})
}
var allusers []bson.M
if err = result.All(ctx, &allusers); err!=nil{
	log.Fatal(err)
}
c.JSON(http.StatusOK, allusers[0])}}

func GetUser() gin.HandlerFunc{
	return func(c *gin.Context){
		userId := c.Param("user_id")

		if err := helper.MatchUserTypeToUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id":userId}).Decode(&user)
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

/*func AlertCon() gin.HandlerFunc {
	return func(c *gin.Context) {
		var alertConfig models.AlertConfig

		if err := c.ShouldBindJSON(&alertConfig); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		filter := bson.M{"id": alertConfig.ID}
		var existingAlert models.AlertConfig
		err := alertCollection.FindOne(context.Background(), filter).Decode(&existingAlert)
		if err == nil {
			update := bson.M{"$set": bson.M{
				"description": alertConfig.Description,
				"severity":    alertConfig.Severity,
			}}
			_, err := alertCollection.UpdateOne(context.Background(), filter, update)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update alert"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Alert updated successfully", "alertConfig": alertConfig})
		} else if err == mongo.ErrNoDocuments {
			_, err := alertCollection.InsertOne(context.Background(), alertConfig)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert alert"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Alert inserted successfully", "alertConfig": alertConfig})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query database"})
		}
	}
}*/

func PostRem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var alert models.Alerts
		//var alertOutput models.AlertOutput

		// Bind alert and alertOutput from the request
		if err := c.ShouldBindJSON(&alert); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		/*if err := c.ShouldBindJSON(&alertOutput); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}*/

		// Get the user ID from the context or request headers
		userIDHandler := middleware.GetUserIdContext()
		userIDHandler(c)
		userID := c.GetString("uid")
		//fmt.Print("userid=", userID)
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user ID not found"})
			return
		}

		
		// Check for duplicate alert
        if isDuplicate, err := checkForDuplicateAlert(alert); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check for duplicate alert"})
            return
        } else if isDuplicate {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Duplicate alert found"})
            return
        }
        

		// Insert the alert into the alerts collection
		result, err := alertCollection.InsertOne(context.Background(), alert)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert alert"})
			return
		}

		// Add the alert ID to the user's Alerts array field

		/*res, err := userCollection.update(
   		{ _id: ObjectId(userID) }, // Query to find the document
   		{ $push: { Alerts: { $each: [alert.ID] } } } // Update operation
		)*/

		alertIDString := alert.ID.String()

		// Define the filter to match the userID
		filter := bson.M{"user_id": userID}

		// Define the update operation to push alertIDString into the Alert array
		update := bson.M{"$push": bson.M{"Alert": alertIDString}}

		// Perform the update operation
		res, err := userCollection.UpdateOne(context.Background(), filter, update)
		if err != nil {
    		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user document"})
    		return
		}

		// Use res for error checking or logging if needed
		if res.ModifiedCount == 0 {
    		// Log a message indicating that no documents were modified
    		log.Println("No documents were modified during the update operation")
		} else {
    		// Log a message indicating the number of documents modified
    		log.Printf("%d document(s) were modified during the update operation\n", res.ModifiedCount)
		}

		// Return a success response with the inserted alert
		c.JSON(http.StatusOK, gin.H{"message": "Alert inserted successfully", "alertID": result.InsertedID})
	}
}

func checkForDuplicateAlert(alert models.Alerts) (bool, error) {
    // Query the MongoDB collection to check for duplicates
    // Exclude _id field from the filter
    filter := bson.M{
        "id":          alert.ID,
        "nodeid":      alert.NodeID,
        "description": alert.Description,
        "severity":    alert.Severity,
        "source":      alert.Source,
        "createdat":   alert.CreatedAt,
        // Add other fields as needed
    }

    // Count the number of documents that match the filter
    count, err := alertCollection.CountDocuments(context.Background(), filter)
    if err != nil {
        // Return error if there's any issue with the database query
        return false, err
    }

    // If count is greater than 0, a duplicate alert exists
    return count > 0, nil
}


func AlertConfig() gin.HandlerFunc {
    return func(c *gin.Context) {
        fmt.Println("IN ALERT CONFIG")
        userID := c.GetString("email")
        if userID == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found"})
            return
        }

		// fmt.println()

        var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()
        var user models.User
		err := userCollection.FindOne(ctx, bson.M{"email": userID}).Decode(&user)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        var alertConfig struct {
            Categories []string `json:"categories"`
            Severities []string `json:"severities"`
        }
        if err := c.BindJSON(&alertConfig); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        filter := bson.M{"email": userID}
        update := bson.M{
            "$set": bson.M{
                "alert.categories": alertConfig.Categories,
                "alert.severities": alertConfig.Severities,
            },
        }
        _, err = userCollection.UpdateOne(ctx, filter, update)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user document"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Alert configuration saved successfully"})
    }
}

