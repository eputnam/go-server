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
	store, err := db.NewStore(globalConfig)
	if nil != err {
		panic(err)
	}
	server := api.NewServer(store)
	server.StartServer(globalConfig.Server.Host + ":" + globalConfig.Server.Port)
}
