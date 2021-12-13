package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type SurveyDB struct {
	gorm.Model
	Team        string
	ResponseURL string
	Active      bool
	EndTime     int64
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

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "crazytownbananapants"
	dbname   = "postgres"
)

func NewDatabase() *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if nil != err {
		panic(err)
	}

	err = db.AutoMigrate(&SurveyDB{})
	handleError(err)
	err = db.AutoMigrate(&ResponseDB{})
	handleError(err)
	err = db.AutoMigrate(&QuestionDB{})
	handleError(err)

	return db
}

func handleError(err error) {
	if nil != err {
		panic(err)
	}
}
