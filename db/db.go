package db

import (
	"fmt"
	"github.com/eputnam/health-check-server/config"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

type DataStore struct {
	DB *gorm.DB
}

type Survey struct {
	gorm.Model
	TeamID      uint
	ResponseURL string
	EndTime     int64
}

type QuestionResponse struct {
	gorm.Model
	SurveyID   uint
	QuestionID uint
	Answer     int
	UserID     string
}

type Question struct {
	gorm.Model
	Text     string
	SurveyID uint
	Order    int
}

type Team struct {
	gorm.Model
	Name string
}

func (team Team) toLogString() string {
	return fmt.Sprintf("Name: %s, ID=%d", team.Name, team.ID)
}

func (Team) TableName() string {
	return "teams"
}

func (ds *DataStore) SaveQuestion(q Question) Question {
	if result := ds.DB.Create(&q); nil != result.Error {
		panic(result.Error)
	}
	logrus.Debugf("Successfully saved question %d", q.ID)
	return q
}

func (ds *DataStore) GetQuestionsForSurvey(surveyId uint64) []Question {
	var questions []Question
	if result := ds.DB.Find(&questions, "survey_id = ?", surveyId); nil != result.Error {
		panic(result.Error)
	}
	logrus.Debugf("Successfully found questions for survey %d", surveyId)
	return questions
}

func (ds *DataStore) SaveSurvey(survey Survey) Survey {
	if result := ds.DB.Create(&survey); nil != result.Error {
		panic(result.Error)
	}
	logrus.Debugf("Successfuly saved survey %d", survey.ID)
	return survey
}

func (ds *DataStore) GetSurvey(id uint64) (Survey, error) {
	var survey Survey
	if result := ds.DB.First(&survey, id); nil != result.Error {
		return Survey{}, result.Error
	}
	logrus.Debugf("Succesfully retrieved survey %d", survey.ID)
	return survey, nil
}

func (ds *DataStore) GetSurveysByTeam(teamId uint) ([]Survey, error) {
	var surveys []Survey
	if result := ds.DB.Find(&surveys, "team_id = ?", teamId); nil != result.Error {
		return []Survey{}, result.Error
	}
	logrus.Debugf("Succesfully retrieved teams for team %d", teamId)
	return surveys, nil
}

func (ds *DataStore) SaveTeam(team Team) Team {
	if result := ds.DB.Create(&team); nil != result.Error {
		panic(result.Error)
	}
	logrus.Debugf("Successfully saved team %s", team.toLogString())
	return team
}

func (ds *DataStore) GetTeams() []Team {
	var teams []Team
	if result := ds.DB.Find(&teams, &Team{}); nil != result.Error {
		panic(result.Error)
	}
	logrus.Debug("Successfully got list of teams")
	return teams
}

func NewStore(conf config.GlobalConfig) (*DataStore, error) {
	dbConf := conf.DB
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbConf.Host, dbConf.Port, dbConf.User, dbConf.Password, dbConf.DBName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger(conf)})
	if nil != err {
		return nil, err
	}
	logrus.Infof("Connected to database at %s:%s", dbConf.Host, dbConf.Port)

	if err := db.AutoMigrate(&Survey{}); nil != err {
		return nil, err
	}
	if err := db.AutoMigrate(&QuestionResponse{}); nil != err {
		return nil, err
	}
	if err := db.AutoMigrate(&Question{}); nil != err {
		return nil, err
	}
	if err := db.AutoMigrate(&Team{}); nil != err {
		return nil, err
	}

	return &DataStore{DB: db}, nil
}

func newLogger(conf config.GlobalConfig) logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{LogLevel: conf.GetDbLogLevel(), Colorful: true})
}
