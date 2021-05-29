package factory

import (
	"log"

	"github.com/TruckX/dao/mongodb"
)

func FactoryDao(e string) {
	switch e {
	case "mongodb":
		mongodb.InitDB()
	default:
		log.Fatalf("This engine %s is not implemented", e)
		return
	}
}
