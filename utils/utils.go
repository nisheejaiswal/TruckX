package utils

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/TruckX/models"
	"github.com/spf13/viper"
)

var Config models.Config

// InitConfig setups configuration
func InitConfig() {

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&Config)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	fmt.Println("Database URL is\t", Config.DatabaseEngine+"://"+Config.DatabaseServer+":"+Config.DatabasePort)
}

// StringInSlice returns true if string exists in slice
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// PrettyJSON returns a JSON indeneted output
func PrettyJSON(body interface{}) error {
	prettyJSON, err := json.MarshalIndent(body, "", "    ")
	if err != nil {
		log.Println("Failed to generate json", err)
		return err
	}
	fmt.Printf("%s\n", string(prettyJSON))
	return nil
}
