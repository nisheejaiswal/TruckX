package services

import "go.mongodb.org/mongo-driver/mongo"

type MongoInstance struct {
	dbClient *mongo.Client
	dbName   string
}
