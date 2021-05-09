package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

//BaseConfig is the base config used in most applications
type BaseConfig struct {
	APIURL    string
	Secret    string
	Database  DatabaseConfig
	Messaging MessagingConfig
}

//DatabaseConfig is the config used to connect to the central database
type DatabaseConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

//MessagingConfig is the config used to connect to the RabbitMQ server
type MessagingConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

var config BaseConfig = BaseConfig{
	APIURL: "http://localhost:8080",
	Secret: "RushmeadLikesGettingHisBackdoorSmashedIn",
	Database: DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		Username: "postgres",
		Password: "",
		Database: "do",
	},
	Messaging: MessagingConfig{
		Host:     "localhost",
		Port:     5672,
		Username: "app",
		Password: "thisIsmyPassword",
	},
}

// LoadConfig Loads the config from the config.json file
func LoadConfig() {
	data, err := ioutil.ReadFile("config.json")
	if err == nil {
		err = json.Unmarshal(data, &config)
	} else {
		jsonBytes, err := json.Marshal(config)
		if err != nil {
			fmt.Println("Whoops unable to JSONify the new config")
		} else {
			err = ioutil.WriteFile("config.json", jsonBytes, 0644)
			if err != nil {
				fmt.Println("An error occured whilst writing the new config to file, but we don't really care cause the software can still run")
			}
		}
	}
}

// GetConfig returns the currently loaded BaseConfig
func GetConfig() BaseConfig {
	return config
}
