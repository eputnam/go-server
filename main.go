package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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

func main() {
	initializeDB()

	router := gin.Default()
	router.Run("localhost:8080")
}
