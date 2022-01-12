package api

import (
	"errors"
	"fmt"
	"github.com/eputnam/health-check-server/db"
	"github.com/eputnam/health-check-server/vcs"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

type QuestionResponseAPI struct {
	QuestionID uint   `json:"questionId"`
	SurveyID   uint   `json:"surveyId"`
	Response   string `json:"response"`
}

type QuestionAPI struct {
	ID        uint                  `json:"id"`
	SurveyID  uint                  `json:"surveyId"`
	Text      string                `json:"text"`
	Order     int                   `json:"order"`
	Responses []QuestionResponseAPI `json:"responses"`
}

type SurveyAPI struct {
	ID           uint          `json:"id"`
	Team         TeamAPI       `json:"team"`
	Questions    []QuestionAPI `json:"questions"`
	Active       bool          `json:"active"`
	MaxResponses int           `json:"maxResponses"`
}

type CreateSurveyRequest struct {
	TeamID       string `json:"teamId"`
	MaxResponses int    `json:"maxResponses"`
}

type GetSurveyRequest struct {
	ID uint64 `uri:"id" binding:"required"`
}

func (server *Server) CreateSurvey(ginc *gin.Context) {
	var request CreateSurveyRequest
	err := ginc.BindJSON(&request)
	logrus.Tracef("Received CreateSurveyRequest: %s", request.toLogString())
	if nil != err {
		ginc.JSON(http.StatusBadRequest, err)
		panic(err)
	}

	surveyTimeoutMilli := int64(6 * 60 * 1000)
	endTime := time.Now().UnixMilli() + surveyTimeoutMilli
	teamId, err := strconv.ParseUint(request.TeamID, 10, 1)
	newSurvey := db.Survey{TeamID: teamId, EndTime: endTime, ResponseURL: "placeholder"}
	surveyDb := server.DB.SaveSurvey(newSurvey)

	client := vcs.NewClient(server.Config.GitHub)
	questions, err := client.GetQuestionsList()
	if nil != err {
		ginc.JSON(http.StatusInternalServerError, err)
		panic(err)
	}

	var qr []QuestionAPI

	for i, v := range questions.Questions {
		questionDB := server.DB.SaveQuestion(db.Question{Text: v, Order: i, SurveyID: surveyDb.ID})
		qr = append(qr, questionDbToApi(questionDB))
	}

	surveyApi := surveyDbToApi(surveyDb, qr)

	ginc.JSON(http.StatusOK, &surveyApi)
}

func (server *Server) GetSurvey(ginc *gin.Context) {
	var request GetSurveyRequest
	err := ginc.ShouldBindUri(&request)
	if nil != err {
		ginc.JSON(http.StatusBadRequest, err)
		return
	}
	surveyDb, err := server.DB.GetSurvey(request.ID)
	logrus.Info(surveyDb)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ginc.JSON(http.StatusNotFound, err)
		return
	}

	questions := server.DB.GetQuestionsForSurvey(request.ID)

	surveyApi := surveyDbToApi(surveyDb, questionsDbToApi(questions))

	ginc.JSON(http.StatusOK, &surveyApi)
}

func (server *Server) GetSurveysByTeam(ginc *gin.Context) {

}

func surveyDbToApi(s db.Survey, q []QuestionAPI) SurveyAPI {
	active := time.Now().UnixMilli() > s.EndTime
	return SurveyAPI{ID: s.ID, Team: TeamAPI{ID: s.Team.ID, Name: s.Team.Name}, Active: active, Questions: q}
}

func questionDbToApi(q db.Question) QuestionAPI {
	return QuestionAPI{ID: q.ID, Text: q.Text, SurveyID: q.SurveyID, Order: q.Order}
}

func questionsDbToApi(qdb []db.Question) []QuestionAPI {
	var qapi []QuestionAPI
	for _, v := range qdb {
		qapi = append(qapi, questionDbToApi(v))
	}
	return qapi
}

func (c *CreateSurveyRequest) toLogString() string {
	return fmt.Sprintf("TeamID=%d MaxResponses=%d", c.TeamID, c.MaxResponses)
}
