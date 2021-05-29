package worker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/TruckX/models"
	"github.com/TruckX/utils"
)

func LocationMessage(d time.Duration) {
	for time := range time.NewTicker(d).C {
		vehiclesList, err := GetVehiclesList()
		if err != nil {
			log.Println(err)
		}
		for _, vehicle := range vehiclesList {
			message, err := SendVehicleLocation(vehicle.IMEI, time)
			if err != nil {
				log.Println(err)
			}
			utils.PrettyJSON(message)
		}
	}
}

func SendVehicleLocation(imei string, time time.Time) (models.LocationMessage, error) {
	url := utils.Config.Proto + utils.Config.DatabaseServer + ":8081" + "/api/v1/vehicle/" + imei
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		return models.LocationMessage{}, err
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error :  %s\n", err)
		return models.LocationMessage{}, err
	}

	var vehicles models.Vehicle
	err = json.Unmarshal(data, &vehicles)
	if err != nil {
		fmt.Printf("Error :  %s\n", err)
		return models.LocationMessage{}, err
	}

	coordinates, err := FindLocation(vehicles.Region)
	if err != nil {
		fmt.Printf("Error :  %s\n", err)
		return models.LocationMessage{}, err
	}

	locationMessage := models.LocationMessage{
		Type:         "LOCATION",
		LocationTime: time,
		Latitude:     coordinates.Latitude,
		Longitude:    coordinates.Longitude,
	}

	return locationMessage, nil

}

func FindLocation(region string) (models.Coordinates, error) {
	url := utils.Config.Proto + utils.Config.DatabaseServer + ":8081" + "/api/v1/location/" + region
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		return models.Coordinates{}, err
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error :  %s\n", err)
		return models.Coordinates{}, err
	}

	var coordinates models.Coordinates
	err = json.Unmarshal(data, &coordinates)
	if err != nil {
		fmt.Printf("Error :  %s\n", err)
		return models.Coordinates{}, err
	}

	return coordinates, nil
}

func LoginMessage(d time.Duration) {
	for time := range time.NewTicker(d).C {
		vehiclesList, err := GetVehiclesList()
		if err != nil {
			log.Println(err)
		}
		for _, vehicle := range vehiclesList {
			if vehicle.PowerOn {
				err := GetLoginMessage(vehicle.IMEI, time)
				if err != nil {
					log.Println(err)
					return
				}
			}
		}
	}
}

func GetLoginMessage(imei string, time time.Time) error {
	url := utils.Config.Proto + utils.Config.DatabaseServer + ":8081" + "/api/v1/vehicle/login/" + imei
	response, err := http.Post(url, "", nil)
	if err != nil {
		log.Printf("The HTTP request failed with error %s\n", err)
	}
	var message models.LoginMessage

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return err
	}

	err = json.Unmarshal(data, &message)
	if err != nil {
		log.Println("Unmarshal Error!")
		return err
	}

	err = utils.PrettyJSON(message)
	if err != nil {
		log.Println("PrettyJSON Error!")
		return err
	}

	return nil
}
