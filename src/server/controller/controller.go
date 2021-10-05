package controller

import (
	"chat/src/db"
	"chat/src/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"log"
)

type Controller struct {
	ConfigYaml model.ConfigYaml
}

// 查询聊天消息记录
func (controller Controller) ShowChatMessage(ctx *gin.Context) {

	// 校验请求参数
	var ShowChatMessageParams model.ShowChatMessageParams
	if err := ctx.ShouldBindBodyWith(&ShowChatMessageParams, binding.JSON); err != nil {
		JSONFail(ctx, IllegalRequestParams, fmt.Sprintf(`%s: %s`, RequestParamsErrMessage, err.Error()))
		log.Println(fmt.Sprintf(`%s: %s`, RequestParamsErrMessage, err.Error()))
		return
	}

	var err error
	chatMessage := make([]*model.ChatMessage, 0)

	switch ShowChatMessageParams.UserType {
	case "Admin":
		// 查询聊天消息记录
		chatMessage, err = db.SelectAllChatMessage()
		if err != nil {
			JSONFail(ctx, OperationDBError, fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
			log.Println(fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
			return
		}
	case "Customer":
		// 校验IP地址
		ip := ctx.ClientIP()
		err = CheckIP(ip)
		if err != nil {
			JSONFail(ctx, IpAnalysisError, err.Error())
			return
		}
		// 查询聊天消息记录
		chatMessage, err = db.SelectChatMessage(ip)
		if err != nil {
			JSONFail(ctx, OperationDBError, fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
			log.Println(fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
			return
		}

		firstReply := new(model.FirstReply)
		if len(chatMessage) == 0 {
			// 查询首次回复设置
			firstReply, err = db.SelectFirstReply()
			if err != nil {
				JSONFail(ctx, OperationDBError, fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
				log.Println(fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
				return
			}
			// 创建首次回复消息
			newChatMessage := new(model.ChatMessage)
			newChatMessage.IP = ip
			newChatMessage.Message = firstReply.Message
			newChatMessage.MessageType = MessageTypeFirst
			err = db.CreateChatMessage(newChatMessage)
			if err != nil {
				JSONFail(ctx, OperationDBError, fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
				log.Println(fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
				return
			}
		}
	default:
		JSONFail(ctx, IllegalRequestParams, fmt.Sprintf(`%s: %s`, RequestParamsErrMessage, ShowChatMessageParams.UserType))
		log.Println(fmt.Sprintf(`%s: %s`, RequestParamsErrMessage, ShowChatMessageParams.UserType))
		return
	}

	// 构建返回包
	result := make(map[string][]model.ChatMessageInfo)
	for _, val := range chatMessage {
		_, ok := result[val.IP]
		if !ok {
			result[val.IP] = make([]model.ChatMessageInfo, 0)
		}
		chatMessageInfo := model.ChatMessageInfo{}
		chatMessageInfo.Message = val.Message
		chatMessageInfo.MessageType = val.MessageType
		chatMessageInfo.CreateTime = val.CreatedAt.Format(TimeFormat)
		chatMessageInfo.OptionInfo = make([]model.OptionInfo, 0)
		if val.MessageType == MessageTypeFirst {
			// 查询首次回复选项信息
			firstReplyOptionMessage, err := db.SelectFirstReplyOptionMessage()
			if err != nil {
				JSONFail(ctx, OperationDBError, fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
				log.Println(fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
				return
			}
			for _, tmp := range firstReplyOptionMessage {
				optionInfo := model.OptionInfo{}
				optionInfo.Option = tmp.Option
				optionInfo.Content = tmp.Content
				chatMessageInfo.OptionInfo = append(chatMessageInfo.OptionInfo, optionInfo)
			}
		}
		result[val.IP] = append(result[val.IP], chatMessageInfo)
	}

	JSONSuccess(ctx, result)
}
