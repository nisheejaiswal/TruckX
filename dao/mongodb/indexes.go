package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/TruckX/models"
	"github.com/TruckX/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

func InitDB() {
	collections := GetDBCollectionNames()
	CreateIndexes(collections)
}

// GetDBCollectionNames ...
func GetDBCollectionNames() (collectionNames []models.Collection) {
	collectionNames = []models.Collection{
		{
			CollectionName: "vehicle",
			Fields: []models.Field{
				{
					FieldName: "imei",
					Unique:    true,
				},
			},
		},
		{
			CollectionName: "location",
			Fields: []models.Field{
				{
					FieldName: "location",
					Unique:    true,
				},
			},
		},
		{
			CollectionName: "video",
			Fields: []models.Field{
				{
					FieldName: "imei",
					Unique:    true,
				},
			},
		},
	}
	return
}

// CreateIndexes - creates an index for a specific field in a collection
func CreateIndexes(collections []models.Collection) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, collection := range collections {
		models := []mongo.IndexModel{}

		for _, field := range collection.Fields {
			models = append(models, mongo.IndexModel{
				Keys:    bson.M{field.FieldName: 1},
				Options: options.Index().SetUnique(field.Unique),
			})
		}

		database := ConnectDB().Database(utils.Config.DatabaseName).Collection(collection.CollectionName)
		_, err := database.Indexes().CreateMany(ctx, models)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}
