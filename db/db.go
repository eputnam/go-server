package db

import (
	"fmt"
	"github.com/eputnam/health-check-server/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DataStore struct {
	DB *gorm.DB
}

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

type TeamDB struct {
	gorm.Model
	Name string
}

func (TeamDB) TableName() string {
	return "teams"
}

func (ds *DataStore) SaveTeam(team TeamDB) TeamDB {
	result := ds.DB.Create(&team)
	if nil != result.Error {
		panic(result.Error)
	}
	return team
}

func (ds *DataStore) GetTeams() []TeamDB {
	var teams []TeamDB
	result := ds.DB.Find(&teams, &TeamDB{})
	if nil != result.Error {
		panic(result.Error)
	}
	return teams
}

func NewStore(conf config.GlobalConfig) *DataStore {
	dbConf := conf.DB
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbConf.Host, dbConf.Port, dbConf.User, dbConf.Password, dbConf.DBName)
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
	err = db.AutoMigrate(&TeamDB{})
	handleError(err)

	return &DataStore{DB: db}
}

func handleError(err error) {
	if nil != err {
		panic(err)
	}
}
