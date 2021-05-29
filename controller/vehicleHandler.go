package controller

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/TruckX/constants"
	"github.com/TruckX/dao/mongodb"
	"github.com/TruckX/dao/services"
	"github.com/TruckX/models"
	"github.com/TruckX/utils"
	jsonpatch "github.com/evanphx/json-patch"
	"github.com/gorilla/mux"
)

func vehicleHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*2)
	defer cancel()

	dbClient := mongodb.ConnectDB()

	var vehicle models.Vehicle
	err := json.NewDecoder(r.Body).Decode(&vehicle)
	if err != nil {
		errorParser(w, "badreq: Invalid Request Body!")
		return
	}

	err = services.VehicleMongoService(dbClient, utils.Config.DatabaseName).CreateVehicle(ctx, vehicle)
	if err != nil {
		errorParser(w, err.Error())
		return
	}

	url := utils.Config.Proto + utils.Config.DatabaseServer + ":8081" + "/api/v1/location/" + vehicle.Region
	response, err := http.Post(url, "", nil)
	if err != nil {
		log.Printf("The HTTP request failed with error %s\n", err)
	}

	_, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}

	url = utils.Config.Proto + utils.Config.DatabaseServer + ":8081" + "/api/v1/video"
	response, err = http.Post(url, "", nil)
	if err != nil {
		log.Printf("The HTTP request failed with error %s\n", err)
	}

	_, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}

	responseWithStaus(w, http.StatusCreated, vehicle)
}

func vehicleListHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*2)
	defer cancel()

	dbClient := mongodb.ConnectDB()

	vehicles, err := services.VehicleMongoService(dbClient, utils.Config.DatabaseName).GetAllVehicle(ctx)
	if err != nil {
		errorParser(w, err.Error())
		return
	}

	responseWithStaus(w, http.StatusOK, vehicles)
}

func vehicleGetHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*2)
	defer cancel()

	imei := mux.Vars(r)["imei"]

	dbClient := mongodb.ConnectDB()
	vehicles, err := services.VehicleMongoService(dbClient, utils.Config.DatabaseName).FindVehicle(ctx, imei)
	if err != nil {
		errorParser(w, err.Error())
		return
	}

	responseWithStaus(w, http.StatusOK, vehicles)
}

func vehiclePowerHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*2)
	defer cancel()

	dbClient := mongodb.ConnectDB()

	imei := mux.Vars(r)["imei"]

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorParser(w, "badreq: Invalid Request Body!")
		return
	}

	var decoder map[string]interface{}
	err = json.Unmarshal(body, &decoder)
	if err != nil {
		errorParser(w, "badreq: Unmarshall Error")
		return
	}

	var vehiclePatch models.Vehicle

	vehiclePatch, err = services.VehicleMongoService(dbClient, utils.Config.DatabaseName).FindVehicle(ctx, imei)
	if err != nil {
		errorParser(w, err.Error())
		return
	}

	vehiclePatchBytes, err := json.Marshal(vehiclePatch)
	if err != nil {
		errorParser(w, "badreq: Marshal Error")
		return
	}

	vehicleDecoder, err := json.Marshal(decoder)
	if err != nil {
		errorParser(w, "badreq: Marshal Error")
		return
	}

	patchedJSON, err := jsonpatch.MergePatch(vehiclePatchBytes, vehicleDecoder)
	if err != nil {
		errorParser(w, "badreq: MergePatch Error")
		return
	}

	err = json.Unmarshal(patchedJSON, &vehiclePatch)
	if err != nil {
		errorParser(w, "badreq: Unmarshal Error")
		return
	}

	err = services.VehicleMongoService(dbClient, utils.Config.DatabaseName).UpdateVehicle(ctx, imei, vehiclePatch)
	if err != nil {
		errorParser(w, err.Error())
		return
	}

	responseWithStaus(w, http.StatusOK, vehiclePatch)
}

func vehicleLoginHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*2)
	defer cancel()

	dbClient := mongodb.ConnectDB()

	imei := mux.Vars(r)["imei"]

	_, err := services.VehicleMongoService(dbClient, utils.Config.DatabaseName).FindVehicle(ctx, imei)
	if err != nil {
		errorParser(w, err.Error())
		return
	}

	_, err = services.VehicleMongoService(dbClient, utils.Config.DatabaseName).FindStatus(ctx, imei)
	if err != nil {
		errorParser(w, err.Error())
		return
	}

	response := models.LoginMessage{
		Type: "LOGIN",
		IMEI: imei,
	}

	responseWithStaus(w, http.StatusOK, response)
}

func vehicleEventHandler(w http.ResponseWriter, r *http.Request) {
	alarmType := mux.Vars(r)["alarm_type"]

	if !utils.StringInSlice(alarmType, constants.AlarmType) {
		return
	}

	imei := mux.Vars(r)["imei"]

	var vehicle models.Vehicle
	url := utils.Config.Proto + utils.Config.DatabaseServer + ":8081" + "/api/v1/vehicle/" + imei
	response, err := http.Get(url)
	if err != nil {
		log.Printf("The HTTP request failed with error %s\n", err)
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal(data, &vehicle)
	if err != nil {
		errorParser(w, "badreq: Unmarshall Error")
		return
	}

	var location models.Location
	url = utils.Config.Proto + utils.Config.DatabaseServer + ":8081" + "/api/v1/location/" + vehicle.Region
	response, err = http.Get(url)
	if err != nil {
		log.Printf("The HTTP request failed with error %s\n", err)
	}
	data, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal(data, &location)
	if err != nil {
		errorParser(w, "badreq: Unmarshall Error")
		return
	}

	var video models.Video

	url = utils.Config.Proto + utils.Config.DatabaseServer + ":8081" + "/api/v1/vehicle/video/" + vehicle.IMEI
	response, err = http.Post(url, "", nil)
	if err != nil {
		log.Printf("The HTTP request failed with error %s\n", err)
	}
	data, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal(data, &video)
	if err != nil {
		errorParser(w, "badreq: Unmarshall Error")
		return
	}

	filesList := []string{}

	alarmMessage := models.AlarmMessage{
		Type:      "ALARM",
		AlarmType: alarmType,
		AlarmTime: time.Now(),
		Latitude:  location.Latitude,
		Longitude: location.Longitude,
		FileList:  append(filesList, video.FileName),
	}

	responseWithStaus(w, http.StatusCreated, alarmMessage)
}
