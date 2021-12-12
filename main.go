package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v40/github"
	"golang.org/x/oauth2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	config_path = "config.yaml"
)

var appConfig AppConfig

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

/*
GITHUB STUFF
*/

const (
	teams_path = "teams.json"
)

type teamsResponse struct {
	Teams []string
}

func githubClient() (*github.Client, context.Context) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: appConfig.GitHub.Token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc), ctx
}

func getTeamsList(client *github.Client, ctx context.Context) teamsResponse {
	byteData, err := getGitHubBytes(client, ctx, teams_path)
	checkError(err)

	var jsonData teamsResponse
	json.Unmarshal(byteData, &jsonData)

	return jsonData
}

func getGitHubBytes(client *github.Client, ctx context.Context, path string) ([]byte, error) {
	content, _, _, err := client.Repositories.GetContents(ctx, appConfig.GitHub.Username, appConfig.GitHub.Repository, path, nil)
	checkError(err)

	byteData, err := base64.StdEncoding.DecodeString(*content.Content)
	checkError(err)

	return byteData, err
}

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
	GitHub struct {
		Username   string `yaml:"username"`
		Repository string `yaml:"repository"`
		Token      string `yaml:"token"`
	} `yaml:"github"`
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
	initializeDB()
	client, ghctx := githubClient()

	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		v1.GET("/teams", func(ginc *gin.Context) {
			ginc.JSON(http.StatusOK, getTeamsList(client, ghctx))
		})
	}

	router.Run("localhost:" + appConfig.Server.Port)
}
