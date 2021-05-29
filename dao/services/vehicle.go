package services

import (
	"context"
	"errors"

	"github.com/TruckX/constants"
	"github.com/TruckX/models"
	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type VehicleDao interface {
	CreateVehicle(ctx context.Context, vehicle models.Vehicle) error
	FindVehicle(ctx context.Context, vehicleID string) (models.Vehicle, error)
	UpdateVehicle(ctx context.Context, imei string, vehicle models.Vehicle) error
	FindStatus(ctx context.Context, vehicleID string) (*models.Vehicle, error)
	GetAllVehicle(ctx context.Context) (*[]models.Vehicle, error)
}

// VehicleMongoService ...
func VehicleMongoService(dbClient *mongo.Client, dataBaseName string) VehicleDao {
	return &MongoInstance{dbClient, dataBaseName}
}

func (mgo *MongoInstance) CreateVehicle(ctx context.Context, vehicle models.Vehicle) error {
	_, err := mgo.dbClient.Database(mgo.dbName).Collection(constants.VehicleCollection).InsertOne(ctx, vehicle)
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

func (mgo *MongoInstance) FindVehicle(ctx context.Context, vehicleID string) (models.Vehicle, error) {
	var result models.Vehicle

	err := mgo.dbClient.Database(mgo.dbName).Collection(constants.VehicleCollection).FindOne(context.TODO(), bson.M{"imei": vehicleID}).Decode(&result)
	if err != nil && err == mongo.ErrNoDocuments {
		return models.Vehicle{}, err
	}

	return result, nil
}

func (mgo *MongoInstance) UpdateVehicle(ctx context.Context, imei string, vehicleBody models.Vehicle) error {
	filter := bson.M{"imei": imei}
	_, err := mgo.dbClient.Database(mgo.dbName).Collection(constants.VehicleCollection).UpdateOne(ctx, filter, bson.M{"$set": vehicleBody})

	return err
}

func (mgo *MongoInstance) FindStatus(ctx context.Context, vehicleID string) (*models.Vehicle, error) {
	var result models.Vehicle

	filter := bson.M{"imei": vehicleID, "power_on": true}
	err := mgo.dbClient.Database(mgo.dbName).Collection(constants.VehicleCollection).FindOne(context.TODO(), filter).Decode(&result)
	if err != nil && err == mongo.ErrNoDocuments {
		return nil, err
	}

	return &result, nil
}

func (mgo *MongoInstance) GetAllVehicle(ctx context.Context) (*[]models.Vehicle, error) {
	var result []models.Vehicle

	cursor, err := mgo.dbClient.Database(mgo.dbName).Collection(constants.VehicleCollection).Find(ctx, bson.M{})
	if err != nil && err == mongo.ErrNoDocuments {
		return nil, err
	}

	err = cursor.All(ctx, &result)
	return &result, err
}
