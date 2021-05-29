package models

type Video struct {
	IMEI     string `json:"imei" bson:"imei"`
	FileName string `json:"file_name" bson:"file_name"`
	Data     string `json:"data" bson:"data"`
}
