package services

import (
	"context"
	"errors"

	"github.com/TruckX/constants"
	"github.com/TruckX/models"
	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type LocationDao interface {
	CreateLocation(ctx context.Context, location models.Location) error
	FindLocation(ctx context.Context, name string) (models.Location, error)
}

// LocationMongoService ...
func LocationMongoService(dbClient *mongo.Client, dataBaseName string) LocationDao {
	return &MongoInstance{dbClient, dataBaseName}
}

func (mgo *MongoInstance) CreateLocation(ctx context.Context, location models.Location) error {
	_, err := mgo.dbClient.Database(mgo.dbName).Collection(constants.LocationCollection).InsertOne(ctx, location)
	if err != nil {
		mgoErr, ok := err.(mongo.WriteException)
		if !ok {
			return errors.New("mongodb error")
		}
		if errCode := mgoErr.WriteErrors[0].Code; errCode == 11000 {
			return errors.New("duplicate: vehicle imei number already exists")
		}
	}

	return nil
}

func (mgo *MongoInstance) FindLocation(ctx context.Context, name string) (models.Location, error) {
	var result models.Location

	err := mgo.dbClient.Database(mgo.dbName).Collection(constants.LocationCollection).FindOne(context.TODO(), bson.M{"name": name}).Decode(&result)
	if err != nil && err == mongo.ErrNoDocuments {
		return models.Location{}, err
	}

	return result, nil
}
