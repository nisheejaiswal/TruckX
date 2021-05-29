package mongodb

import (
	"context"
	"log"

	"github.com/TruckX/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	clientOptions := options.Client().ApplyURI(utils.Config.DatabaseEngine + "://" + utils.Config.DatabaseServer + ":" + utils.Config.DatabasePort)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	return client
}
