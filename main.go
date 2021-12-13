package main

import (
	"os"

	"github.com/eputnam/health-check-server/api"
	"github.com/eputnam/health-check-server/db"
	"gopkg.in/yaml.v2"
)

const (
	config_path = "config.yaml"
)

var appConfig AppConfig

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

/*
CONFIG STUFF
*/
type AppConfig struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
}

func loadConfig() AppConfig {
	file, err := os.Open(config_path)
	checkError(err)
	defer file.Close()

	var config AppConfig
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	checkError(err)

	return config
}

func init() {
	appConfig = loadConfig()
}

func main() {
	db := db.NewDatabase()
	server := api.NewServer(db)
	server.StartServer("localhost:" + appConfig.Server.Port)
}
