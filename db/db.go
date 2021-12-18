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
	if result := ds.DB.Create(&team); nil != result.Error {
		panic(result.Error)
	}
	return team
}

func (ds *DataStore) GetTeams() []TeamDB {
	var teams []TeamDB
	if result := ds.DB.Find(&teams, &TeamDB{}); nil != result.Error {
		panic(result.Error)
	}
	return teams
}

func NewStore(conf config.GlobalConfig) (*DataStore, error) {
	dbConf := conf.DB
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbConf.Host, dbConf.Port, dbConf.User, dbConf.Password, dbConf.DBName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if nil != err {
		return nil, err
	}

	if err := db.AutoMigrate(&SurveyDB{}); nil != err {
		return nil, err
	}
	if err := db.AutoMigrate(&ResponseDB{}); nil != err {
		return nil, err
	}
	if err := db.AutoMigrate(&QuestionDB{}); nil != err {
		return nil, err
	}
	if err := db.AutoMigrate(&TeamDB{}); nil != err {
		return nil, err
	}

	return &DataStore{DB: db}, nil
}
