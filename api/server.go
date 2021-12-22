package api

import (
	"github.com/eputnam/health-check-server/config"
	"github.com/eputnam/health-check-server/db"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/toorop/gin-logrus"
)

type Server struct {
	DB     *db.DataStore
	Router *gin.Engine
	Config config.GlobalConfig
}

func NewServer(db *db.DataStore, c config.GlobalConfig) *Server {
	server := &Server{DB: db, Config: c}

	log := logrus.New()

	router := gin.New()
	router.Use(ginlogrus.Logger(log), gin.Recovery())

	// router stuff
	basePath := "/api"
	v1Group := router.Group(basePath + "/v1")
	{
		v1Group.POST("/teams", server.CreateTeam)
		v1Group.GET("/teams", server.GetTeams)

		v1Group.POST("/surveys", server.CreateSurvey)
		v1Group.GET("/surveys/:id", server.GetSurvey)
	}

	server.Router = router
	return server
}

func (server *Server) StartServer(address string) {
	server.Router.Run(address)
}
