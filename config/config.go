package config

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm/logger"
	"os"
	"strings"
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
		LogLevel string `yaml:"loglevel"`
	} `yaml:"db"`
	Logrus struct {
		Level string `yaml:"level"`
	} `yaml:"logrus"`
}

func LoadConfig() GlobalConfig {
	file, err := os.Open(config_path)
	checkError(err)
	defer file.Close()

	var config GlobalConfig
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); nil != err {
		panic(err)
	}
	logLevel := config.GetLogrusLevel()
	logrus.SetLevel(logLevel)
	logrus.Print("Set logrus level to " + logLevel.String())
	logrus.Infof("Successfully loaded app configuration from %s", config_path)

	return config
}

func checkError(err error) {
	if nil != err {
		panic(err)
	}
}

func (gc *GlobalConfig) GetDbLogLevel() logger.LogLevel {
	configLevel := strings.ToLower(gc.DB.LogLevel)
	switch configLevel {
	case "info":
		return logger.Info
	case "error":
		return logger.Error
	}
	return logger.Silent
}

func (gc *GlobalConfig) GetLogrusLevel() logrus.Level {
	configLevel := strings.ToLower(gc.Logrus.Level)
	switch configLevel {
	case "debug":
		return logrus.DebugLevel
	case "error":
		return logrus.ErrorLevel
	}
	return logrus.InfoLevel
}
