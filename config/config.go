package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

const config_path = "config.yaml"

type GlobalConfig struct {
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"server"`
	DB struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
	}
}

func LoadConfig() GlobalConfig {
	file, err := os.Open(config_path)
	checkError(err)
	defer file.Close()

	var config GlobalConfig
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	checkError(err)

	return config
}

func checkError(err error) {
	if nil != err {
		panic(err)
	}
}
