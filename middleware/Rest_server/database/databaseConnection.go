package database

import(
"fmt"
"log"
"time"
"os"
"context"
"github.com/joho/godotenv"
"go.mongodb.org/mongo-driver/mongo"
"go.mongodb.org/mongo-driver/mongo/options"
)

func DBinstance() *mongo.Client{
	err := godotenv.Load(".env")
	if err!=nil{
		log.Fatal("Error loading .env file")
	}

	MongoDb := os.Getenv("MONGODB_URL")

	client, err:= mongo.NewClient(options.Client().ApplyURI(MongoDb))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	return client
}

var Client *mongo.Client = DBinstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection{
	var collection *mongo.Collection = client.Database("cluster0").Collection(collectionName)
	return collection
}


// FindAlertConfigByID retrieves an alert configuration by its ID from the database.
/*func FindAlertConfigByID(id uuid.UUID) (*alerts.AlertsConfig, error) {
    collection := OpenCollection(Client, AlertsConfigCollectionName)

    var alertConfig alerts.AlertsConfig
    filter := bson.M{"_id": id}

    err := collection.FindOne(context.TODO(), filter).Decode(&alertConfig)
    if err != nil {
        log.Println("Error finding alert configuration by ID:", err)
        return nil, err
    }

    return &alertConfig, nil
}

// UpdateAlertConfig updates an existing alert configuration in the database.
func UpdateAlertConfig(updatedConfig alerts.AlertsConfig) error {
    collection := OpenCollection(Client, AlertsConfigCollectionName)

    filter := bson.M{"_id": updatedConfig.ID}
    update := bson.M{"$set": updatedConfig}

    _, err := collection.UpdateOne(context.TODO(), filter, update)
    if err != nil {
        log.Println("Error updating alert configuration:", err)
        return err
    }

    return nil
}

// SaveAlertConfig saves a new alert configuration to the database.
func SaveAlertConfig(alertConfig *alerts.AlertsConfig) error {
    collection := OpenCollection(Client, AlertsConfigCollectionName)

    _, err := collection.InsertOne(context.TODO(), alertConfig)
    if err != nil {
        log.Println("Error saving alert configuration:", err)
        return err
    }

    return nil
}
*/
