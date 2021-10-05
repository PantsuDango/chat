package routes

import (
	model2 "chat/src/model"
	controller2 "chat/src/server/controller"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Init(ConfigYaml model2.ConfigYaml) *gin.Engine {

	router := gin.Default()
	router.Use(cors.Default())

	Controller := new(controller2.Controller)
	Controller.ConfigYaml = ConfigYaml
	router.POST("/Chat/ShowChatMessage", Controller.ShowChatMessage)

	return router
}
