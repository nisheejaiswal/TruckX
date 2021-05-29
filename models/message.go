package models

import "time"

type LoginMessage struct {
	Type string `json:"type" bson:"type"`
	IMEI string `json:"imei" bson:"imei"`
}

type AlarmMessage struct {
	Type      string
	AlarmType string
	AlarmTime time.Time
	Latitude  float64
	Longitude float64
	FileList  []string
}

type LocationMessage struct {
	Type         string    `json:"type" bson:"type"`
	LocationTime time.Time `json:"location_time" bson:"location_time"`
	Latitude     float64   `json:"latitude" bson:"latitude"`
	Longitude    float64   `json:"longitude" bson:"longitude"`
}
