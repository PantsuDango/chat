package routes

import (
	"chat/src/model"
	"chat/src/server/controller"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Init(ConfigYaml model.ConfigYaml) *gin.Engine {

	router := gin.Default()
	router.Use(cors.Default())

	Controller := new(controller.Controller)
	Controller.ConfigYaml = ConfigYaml
	router.POST("/Chat/ShowChatMessage", Controller.ShowChatMessage)
	router.POST("/Chat/ShowChatIPList", Controller.ShowChatIPList)
	router.POST("/Chat/UpdateIpContentMap", Controller.UpdateIpContentMap)
	router.POST("/Chat/SelectIpContentMap", Controller.SelectIpContentMap)
	router.POST("/Chat/SendChatMessage", Controller.SendChatMessage)
	router.POST("/Chat/AddKeywordRule", Controller.AddKeywordRule)
	router.POST("/Chat/UpdateKeywordRule", Controller.UpdateKeywordRule)
	router.POST("/Chat/ShowKeywordRule", Controller.ShowKeywordRule)

	return router
}
