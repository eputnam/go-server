package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	DB     *gorm.DB
	Router *gin.Engine
}

func NewServer(db *gorm.DB) *Server {
	server := &Server{DB: db}
	router := gin.Default()

	// router stuff

	server.Router = router
	return server
}

func (server *Server) StartServer(address string) {
	server.Router.Run(address)
}
