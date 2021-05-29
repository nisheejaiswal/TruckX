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

func videoHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*2)
	defer cancel()

	dbClient := mongodb.ConnectDB()
	var video = data.Video

	imei := mux.Vars(r)["imei"]
	video.IMEI = imei

	err := services.VideoMongoService(dbClient, utils.Config.DatabaseName).CreateVideo(ctx, video)
	if err != nil {
		errorParser(w, err.Error())
		return
	}

	responseWithStaus(w, http.StatusCreated, video)
}

func videoGetHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*2)
	defer cancel()

	dbClient := mongodb.ConnectDB()

	imei := mux.Vars(r)["imei"]

	video, err := services.VideoMongoService(dbClient, utils.Config.DatabaseName).FindVideo(ctx, imei)
	if err != nil {
		errorParser(w, err.Error())
		return
	}

	responseWithStaus(w, http.StatusOK, video)
}
