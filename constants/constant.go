package constants

import (
	"net/http"

	"github.com/TruckX/models"
)

const (
	VehicleCollection  = "vehicle"
	LocationCollection = "location"
	VideoCollection    = "video"
)

var AlarmType = []string{"vibration", "overspeed", "crash", "hard_accelerated", "hard_brake", "sharp_turn"}

var ErrorCategory = []models.ErrorResponse{
	{
		Error:      "badreq:",
		StatusCode: http.StatusBadRequest,
	},
	{
		Error:      "notfound:",
		StatusCode: http.StatusNotFound,
	},
	{
		Error:      "duplicate:",
		StatusCode: http.StatusConflict,
	},
	{
		Error:      "Error:",
		StatusCode: http.StatusBadRequest,
	},
}
