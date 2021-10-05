package main

import (
	"chat/src/common"
	"chat/src/db"
	"chat/src/model"
	"chat/src/server/routes"
	"log"
)

var ConfigYaml model.ConfigYaml

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ConfigYaml = common.ReadConfig()

	db.OpenDB(ConfigYaml)
	defer db.CloseDB()

	router := routes.Init(ConfigYaml)
	err := router.Run(":" + ConfigYaml.Server.Port)
	if err != nil {
		log.Fatalf("Server run fail: %s", err)
	}
}
