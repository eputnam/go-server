package api

import (
	"github.com/eputnam/health-check-server/db"
	"github.com/gin-gonic/gin"
)

type Server struct {
	DB     *db.DataStore
	Router *gin.Engine
}

func NewServer(db *db.DataStore) *Server {
	server := &Server{DB: db}
	router := gin.Default()

	// router stuff
	basePath := "/api"
	v1Group := router.Group(basePath + "/v1")
	{
		v1Group.POST("/teams", server.CreateTeam)
		v1Group.GET("/teams", server.GetTeams)
	}

	server.Router = router
	return server
}

func (server *Server) StartServer(address string) {
	server.Router.Run(address)
}
