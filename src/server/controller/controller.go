package controller

import (
	db2 "chat/src/db"
	model2 "chat/src/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"log"
)

type Controller struct {
	ConfigYaml model2.ConfigYaml
}

// 查询聊天消息记录
func (controller Controller) ShowChatMessage(ctx *gin.Context) {

	// 校验请求参数
	var ShowChatMessageParams model2.ShowChatMessageParams
	if err := ctx.ShouldBindBodyWith(&ShowChatMessageParams, binding.JSON); err != nil {
		JSONFail(ctx, IllegalRequestParams, fmt.Sprintf(`%s: %s`, RequestParamsErrMessage, err.Error()))
		log.Println(fmt.Sprintf(`%s: %s`, RequestParamsErrMessage, err.Error()))
		return
	}

	var err error
	chatMessage := make([]*model2.ChatMessage, 0)

	switch ShowChatMessageParams.UserType {
	case "Admin":
		// 查询聊天消息记录
		chatMessage, err = db2.SelectAllChatMessage()
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
		chatMessage, err = db2.SelectChatMessage(ip)
		if err != nil {
			JSONFail(ctx, OperationDBError, fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
			log.Println(fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
			return
		}

		firstReply := new(model2.FirstReply)
		if len(chatMessage) == 0 {
			// 查询首次回复设置
			firstReply, err = db2.SelectFirstReply()
			if err != nil {
				JSONFail(ctx, OperationDBError, fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
				log.Println(fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
				return
			}
			// 创建首次回复消息
			newChatMessage := new(model2.ChatMessage)
			newChatMessage.IP = ip
			newChatMessage.Message = firstReply.Message
			newChatMessage.MessageType = MessageTypeFirst
			err = db2.CreateChatMessage(newChatMessage)
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
	result := make(map[string][]model2.ChatMessageInfo)
	for _, val := range chatMessage {
		_, ok := result[val.IP]
		if !ok {
			result[val.IP] = make([]model2.ChatMessageInfo, 0)
		}
		chatMessageInfo := model2.ChatMessageInfo{}
		chatMessageInfo.Message = val.Message
		chatMessageInfo.MessageType = val.MessageType
		chatMessageInfo.CreateTime = val.CreatedAt.Format(TimeFormat)
		chatMessageInfo.OptionInfo = make([]model2.OptionInfo, 0)
		if val.MessageType == MessageTypeFirst {
			// 查询首次回复选项信息
			firstReplyOptionMessage, err := db2.SelectFirstReplyOptionMessage()
			if err != nil {
				JSONFail(ctx, OperationDBError, fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
				log.Println(fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
				return
			}
			for _, tmp := range firstReplyOptionMessage {
				optionInfo := model2.OptionInfo{}
				optionInfo.Option = tmp.Option
				optionInfo.Content = tmp.Content
				chatMessageInfo.OptionInfo = append(chatMessageInfo.OptionInfo, optionInfo)
			}
		}
		result[val.IP] = append(result[val.IP], chatMessageInfo)
	}

	JSONSuccess(ctx, result)
}
