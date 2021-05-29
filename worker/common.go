package worker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/TruckX/models"
	"github.com/TruckX/utils"
)

func GetVehiclesList() ([]models.Vehicle, error) {
	url := utils.Config.Proto + utils.Config.DatabaseServer + ":8081" + "/api/v1/vehicle/all"
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error :  %s\n", err)
		return nil, err
	}

	var vehicles []models.Vehicle
	err = json.Unmarshal(data, &vehicles)
	if err != nil {
		fmt.Printf("Error :  %s\n", err)
		return nil, err
	}

	return vehicles, nil
}
