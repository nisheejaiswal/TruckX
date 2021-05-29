package models

type Vehicle struct {
	IMEI     string `json:"imei" bson:"imei"`
	PowerOn  bool   `json:"power_on" bson:"power_on"`
	PowerOff bool   `json:"power_off" bson:"power_off"`
	Region   string `json:"region" bson:"region"`
}
