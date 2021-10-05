package main

import (
	common2 "chat/src/common"
	db2 "chat/src/db"
	model2 "chat/src/model"
	routes2 "chat/src/server/routes"
	"log"
)

var ConfigYaml model2.ConfigYaml

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ConfigYaml = common2.ReadConfig()

	db2.OpenDB(ConfigYaml)
	defer db2.CloseDB()

	router := routes2.Init(ConfigYaml)
	err := router.Run(":" + ConfigYaml.Server.Port)
	if err != nil {
		log.Fatalf("Server run fail: %s", err)
	}
}
