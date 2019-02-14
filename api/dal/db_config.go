package dal

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Server    Server    `json:"server"`
	Cockroach Cockroach `json:"cockroach"`
}
type Cockroach struct {
	Host   string `json:"host"`
	Port   string `json:"port"`
	User   string `json:"user"`
	DbName string `json:"dbname"`
}
type Server struct {
	Port string `json:"port"`
}

var config Config
var isLoaded = false

func LoadConfiguration(file string) (Config, error) {
	if isLoaded == true {
		return config, nil
	}
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	json.NewDecoder(configFile).Decode(&config)
	return config, err
}
