package main

import (
	"fmt"
	"os"
	"time"

	"github.com/eputnam/health-check-server/api"
	"gopkg.in/yaml.v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	config_path = "config.yaml"
)

var appConfig AppConfig
var db *gorm.DB

/*
DB STUFF
*/

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "crazytownbananapants"
	dbname   = "postgres"
)

type SurveyDB struct {
	gorm.Model
	Team        string
	ResponseURL string
	Active      bool
	EndTime     time.Time
}

func (SurveyDB) TableName() string {
	return "surveys"
}

type ResponseDB struct {
	gorm.Model
	SurveyID   uint
	QuestionID uint
	Answer     int
	UserID     string
}

func (ResponseDB) TableName() string {
	return "responses"
}

type QuestionDB struct {
	gorm.Model
	Text     string
	SurveyID uint
}

func (QuestionDB) TableName() string {
	return "questions"
}

func initializeDB() *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	checkError(err)

	db.AutoMigrate(&SurveyDB{})
	db.AutoMigrate(&ResponseDB{})
	db.AutoMigrate(&QuestionDB{})

	return db
}

/*

 */

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
	db = initializeDB()
}

func main() {
	server := api.NewServer(db)
	server.StartServer("localhost:" + appConfig.Server.Port)
}
