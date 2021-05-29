package main

import (
	"time"

	"github.com/TruckX/controller"
	"github.com/TruckX/dao/factory"
	"github.com/TruckX/utils"
	"github.com/TruckX/worker"
)

func main() {
	utils.InitConfig()
	factory.FactoryDao(utils.Config.DatabaseEngine)
	go worker.LocationMessage(1 * time.Minute)
	go worker.LoginMessage(1 * time.Minute)
	controller.RunController()
}
