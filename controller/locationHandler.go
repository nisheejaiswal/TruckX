package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/TruckX/dao/mongodb"
	"github.com/TruckX/dao/services"
	"github.com/TruckX/data"
	"github.com/TruckX/utils"
	"github.com/gorilla/mux"
)

func locationHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*2)
	defer cancel()

	dbClient := mongodb.ConnectDB()
	name := mux.Vars(r)["name"]

	var location = data.Locations

	for _, loc := range location {
		if loc.Name == name {
			err := services.LocationMongoService(dbClient, utils.Config.DatabaseName).CreateLocation(ctx, loc)
			if err != nil {
				errorParser(w, err.Error())
				return
			}
			responseWithStaus(w, http.StatusCreated, location)
		}
	}
}

func locationGetHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*2)
	defer cancel()

	dbClient := mongodb.ConnectDB()

	name := mux.Vars(r)["name"]

	location, err := services.LocationMongoService(dbClient, utils.Config.DatabaseName).FindLocation(ctx, name)
	if err != nil {
		errorParser(w, err.Error())
		return
	}

	responseWithStaus(w, http.StatusCreated, location)
}
