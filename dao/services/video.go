package services

import (
	"context"
	"errors"

	"github.com/TruckX/constants"
	"github.com/TruckX/models"
	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type VideoDao interface {
	CreateVideo(ctx context.Context, video models.Video) error
	FindVideo(ctx context.Context, imei string) (models.Video, error)
}

// LocationMongoService ...
func VideoMongoService(dbClient *mongo.Client, dataBaseName string) VideoDao {
	return &MongoInstance{dbClient, dataBaseName}
}

func (mgo *MongoInstance) CreateVideo(ctx context.Context, video models.Video) error {
	_, err := mgo.dbClient.Database(mgo.dbName).Collection(constants.VideoCollection).InsertOne(ctx, video)
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

func (mgo *MongoInstance) FindVideo(ctx context.Context, imei string) (models.Video, error) {
	var result models.Video

	err := mgo.dbClient.Database(mgo.dbName).Collection(constants.VideoCollection).FindOne(context.TODO(), bson.M{"imei": imei}).Decode(&result)
	if err != nil && err == mongo.ErrNoDocuments {
		return models.Video{}, err
	}

	return result, nil
}
