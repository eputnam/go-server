package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type answer struct {
	QuestionID string `json:"questionId"`
	Answer     int    `json:"answer"`
}

type surveyResponse struct {
	ID      string   `json:"id"`
	Answers []answer `json:"answers"`
}

type survey struct {
	ID        string           `json:"id"`
	Team      string           `json:"team"`
	Responses []surveyResponse `json:"responses"`
}

var surveys = []survey{
	{ID: "1", Team: "CD4PE", Responses: nil},
	{ID: "2", Team: "Froyo", Responses: nil},
	{ID: "3", Team: "Cygnus", Responses: nil},
}

/*
    POST /surveys
    	creates a survey for cd4pe
		creates a survey URI

	POST /surveys/response?survey={surveyId}
		creates a survey response
*/

func main() {

	router := gin.Default()
	router.GET("/surveys", getSurveys)
	router.POST("/surveys", postSurvey)

	router.Run("localhost:8080")
}

func getSurveys(c *gin.Context) {
	if c.Query("team") != "" {
		team := c.Query("team")

		for _, s := range surveys {
			if s.Team == team {
				c.IndentedJSON(http.StatusOK, s)
				return
			}
		}
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "survey not found"})
	} else {
		c.IndentedJSON(http.StatusOK, surveys)
	}
}

func postSurvey(c *gin.Context) {
	var newSurvey survey

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newSurvey); err != nil {
		return
	}

	// Add the new album to the slice.
	surveys = append(surveys, newSurvey)
	c.IndentedJSON(http.StatusCreated, newSurvey)
}

func getSurveyByTeam(c *gin.Context) {
	team := c.Query("team")

	for _, s := range surveys {
		if s.Team == team {
			c.IndentedJSON(http.StatusOK, s)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "survey not found"})
}
