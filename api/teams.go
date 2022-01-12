package api

import (
	"github.com/eputnam/health-check-server/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TeamAPI struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type CreateTeamRequest struct {
	Name string `json:"name"`
}

func (server *Server) CreateTeam(ginc *gin.Context) {
	var request CreateTeamRequest
	err := ginc.BindJSON(&request)
	if nil != err {
		ginc.JSON(http.StatusBadRequest, err)
		panic(err)
	}

	newTeam := db.Team{Name: request.Name}

	team := server.DB.SaveTeam(newTeam)

	ginc.JSON(http.StatusOK, &TeamAPI{Name: team.Name, ID: team.ID})
}

func (server *Server) GetTeams(ginc *gin.Context) {
	teams := server.DB.GetTeams()
	var apiTeams []TeamAPI
	for _, v := range teams {
		apiTeams = append(apiTeams, TeamAPI{Name: v.Name, ID: v.ID})
	}
	ginc.JSON(http.StatusOK, &apiTeams)
}
