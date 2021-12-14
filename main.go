package main

import (
	"github.com/eputnam/health-check-server/api"
	"github.com/eputnam/health-check-server/config"
	"github.com/eputnam/health-check-server/db"
)

var globalConfig config.GlobalConfig

func init() {
	globalConfig = config.LoadConfig()
}

func main() {
	store := db.NewStore(globalConfig)
	server := api.NewServer(store)
	server.StartServer(globalConfig.Server.Host + ":" + globalConfig.Server.Port)
}
