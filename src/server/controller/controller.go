package controller

import (
	"chat/src/db"
	"chat/src/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jinzhu/gorm"
	"log"
	"strings"
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
		// 校验IP地址
		ip := ShowChatMessageParams.IP
		if len(ip) != 0 {
			err = CheckIP(ip)
			if err != nil {
				JSONFail(ctx, IpAnalysisError, err.Error())
				return
			}
			// 查询聊天消息记录
			chatMessage, err = db.SelectChatMessage(ip)
		} else {
			// 查询聊天消息记录
			chatMessage, err = db.SelectAllChatMessage()
		}
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

// 查询聊天消息IP列表
func (controller Controller) ShowChatIPList(ctx *gin.Context) {

	// 查询聊天消息记录
	chatMessage, err := db.SelectAllChatMessageIP()
	if err != nil {
		JSONFail(ctx, OperationDBError, fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
		log.Println(fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
		return
	}

	// 构建返回包
	tmpMap := make(map[string]string)
	result := make([]*model.ChatMessageIPListInfo, 0)
	for _, val := range chatMessage {
		if _, ok := tmpMap[val.IP]; !ok {
			tmpMap[val.IP] = val.IP
		} else {
			continue
		}
		chatMessageInfo := new(model.ChatMessageIPListInfo)
		chatMessageInfo.IP = val.IP
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
		// 查询IP备注
		ipContentMap, err := db.SelectIpContentMapByIP(val.IP)
		if err != nil && err != gorm.ErrRecordNotFound {
			JSONFail(ctx, OperationDBError, fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
			log.Println(fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
			return
		} else if err == nil {
			chatMessageInfo.IPContent = ipContentMap.Content
		}
		result = append(result, chatMessageInfo)
	}

	JSONSuccess(ctx, result)
}

// 修改IP备注
func (controller Controller) UpdateIpContentMap(ctx *gin.Context) {

	// 校验请求参数
	var UpdateIpContentMapParams model.UpdateIpContentMapParams
	if err := ctx.ShouldBindBodyWith(&UpdateIpContentMapParams, binding.JSON); err != nil {
		JSONFail(ctx, IllegalRequestParams, fmt.Sprintf(`%s: %s`, RequestParamsErrMessage, err.Error()))
		log.Println(fmt.Sprintf(`%s: %s`, RequestParamsErrMessage, err.Error()))
		return
	}

	// 校验IP地址
	err := CheckIP(UpdateIpContentMapParams.IP)
	if err != nil {
		JSONFail(ctx, IpAnalysisError, err.Error())
		return
	}

	// 查询IP备注
	ipContentMap, err := db.SelectIpContentMapByIP(UpdateIpContentMapParams.IP)
	if err != nil && err != gorm.ErrRecordNotFound {
		JSONFail(ctx, OperationDBError, fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
		log.Println(fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
		return
	} else if err == gorm.ErrRecordNotFound {
		ipContentMap.IP = UpdateIpContentMapParams.IP
	}

	// 更新IP备注
	ipContentMap.Content = UpdateIpContentMapParams.IpContent
	err = db.SaveIpContentMap(ipContentMap)
	if err != nil {
		JSONFail(ctx, OperationDBError, fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
		log.Println(fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
		return
	}

	JSONSuccess(ctx, SuccessMessage)
}

// 查询IP备注
func (controller Controller) SelectIpContentMap(ctx *gin.Context) {

	ipContentMap, err := db.SelectIpContentMap()
	if err != nil {
		JSONFail(ctx, OperationDBError, fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
		log.Println(fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
		return
	}

	JSONSuccess(ctx, ipContentMap)
}

// 发送聊天消息
func (controller Controller) SendChatMessage(ctx *gin.Context) {

	// 校验请求参数
	var SendChatMessageParams model.SendChatMessageParams
	if err := ctx.ShouldBindBodyWith(&SendChatMessageParams, binding.JSON); err != nil {
		JSONFail(ctx, IllegalRequestParams, fmt.Sprintf(`%s: %s`, RequestParamsErrMessage, err.Error()))
		log.Println(fmt.Sprintf(`%s: %s`, RequestParamsErrMessage, err.Error()))
		return
	}

	switch SendChatMessageParams.MessageType {
	case MessageTypeOption, MessageTypeManual, MessageTypeCustomer:
	default:
		JSONFail(ctx, IllegalRequestParams, fmt.Sprintf(`%s`, RequestParamsErrMessage))
		log.Println(fmt.Sprintf(`%s`, RequestParamsErrMessage))
		return
	}

	chatMessage := new(model.ChatMessage)
	chatMessage.Message = SendChatMessageParams.Message
	chatMessage.MessageType = SendChatMessageParams.MessageType
	if SendChatMessageParams.MessageType == MessageTypeCustomer {
		chatMessage.IP = ctx.ClientIP()
	} else {
		chatMessage.IP = SendChatMessageParams.IP
	}
	// 校验IP地址
	err := CheckIP(chatMessage.IP)
	if err != nil {
		JSONFail(ctx, IpAnalysisError, err.Error())
		return
	}
	// 保存聊天消息
	err = db.CreateChatMessage(chatMessage)
	if err != nil {
		JSONFail(ctx, OperationDBError, fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
		log.Println(fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
		return
	}

	if chatMessage.MessageType == MessageTypeCustomer {
		// 查询关键词规则
		keywordRule, err := db.SelectKeywordRuleBySwitch()
		if err != nil {
			JSONFail(ctx, OperationDBError, fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
			log.Println(fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
			return
		}

		keywordMap := make(map[string]string, 0)
		for _, val := range keywordRule {
			// 查询关键词
			keywordRuleMap, err := db.SelectKeywordRuleMapByRuleName(val.RuleName)
			if err != nil {
				JSONFail(ctx, OperationDBError, fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
				log.Println(fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
				return
			}
			for _, tmp := range keywordRuleMap {
				keywordMap[tmp.Content] = val.Content
			}
		}

		// 关键词匹配
		for keyword, content := range keywordMap {
			if !strings.Contains(SendChatMessageParams.Message, keyword) {
				continue
			}
			autoChatMessage := new(model.ChatMessage)
			autoChatMessage.IP = chatMessage.IP
			autoChatMessage.Message = content
			autoChatMessage.MessageType = MessageTypeKeyword
			// 保存聊天消息
			err = db.CreateChatMessage(autoChatMessage)
			if err != nil {
				JSONFail(ctx, OperationDBError, fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
				log.Println(fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
				return
			}
			break
		}
	}

	JSONSuccess(ctx, SuccessMessage)
}

// 添加关键词规则
func (controller Controller) AddKeywordRule(ctx *gin.Context) {

	// 校验请求参数
	var AddKeywordRuleParams model.AddKeywordRuleParams
	if err := ctx.ShouldBindBodyWith(&AddKeywordRuleParams, binding.JSON); err != nil || len(AddKeywordRuleParams.Keyword) == 0 {
		JSONFail(ctx, IllegalRequestParams, fmt.Sprintf(`%s: %s`, RequestParamsErrMessage, err.Error()))
		log.Println(fmt.Sprintf(`%s: %s`, RequestParamsErrMessage, err.Error()))
		return
	}

	keywordRule := new(model.KeywordRule)
	keywordRule.RuleName = AddKeywordRuleParams.RuleName
	keywordRule.Switch = 0
	keywordRule.Content = AddKeywordRuleParams.Content
	// 创建关键词规则
	err := db.CreateKeywordRule(keywordRule)
	if err != nil {
		JSONFail(ctx, OperationDBError, fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
		log.Println(fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
		return
	}

	for _, val := range AddKeywordRuleParams.Keyword {
		keywordRuleMap := new(model.KeywordRuleMap)
		keywordRuleMap.RuleName = AddKeywordRuleParams.RuleName
		keywordRuleMap.Content = val
		// 创建关键词规则
		err = db.CreateKeywordRuleMap(keywordRuleMap)
		if err != nil {
			JSONFail(ctx, OperationDBError, fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
			log.Println(fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
			return
		}
	}

	JSONSuccess(ctx, SuccessMessage)
}

// 更新关键词规则
func (controller Controller) UpdateKeywordRule(ctx *gin.Context) {

	// 校验请求参数
	var UpdateKeywordRuleParams model.UpdateKeywordRuleParams
	if err := ctx.ShouldBindBodyWith(&UpdateKeywordRuleParams, binding.JSON); err != nil {
		JSONFail(ctx, IllegalRequestParams, fmt.Sprintf(`%s: %s`, RequestParamsErrMessage, err.Error()))
		log.Println(fmt.Sprintf(`%s: %s`, RequestParamsErrMessage, err.Error()))
		return
	}

	// 查询关键词规则
	keywordRule, err := db.SelectKeywordRuleByID(UpdateKeywordRuleParams.RuleID)
	if err != nil {
		JSONFail(ctx, OperationDBError, fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
		log.Println(fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
		return
	}

	keywordRule.RuleName = UpdateKeywordRuleParams.RuleName
	keywordRule.Switch = UpdateKeywordRuleParams.Switch
	keywordRule.Content = UpdateKeywordRuleParams.Content
	// 更新关键词规则
	err = db.SaveKeywordRule(keywordRule)
	if err != nil {
		JSONFail(ctx, OperationDBError, fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
		log.Println(fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
		return
	}

	// 删除规则与关键词映射
	err = db.DeleteKeywordRuleMapByRuleName(UpdateKeywordRuleParams.RuleName)
	if err != nil {
		JSONFail(ctx, OperationDBError, fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
		log.Println(fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
		return
	}

	for _, val := range UpdateKeywordRuleParams.Keyword {
		keywordRuleMap := new(model.KeywordRuleMap)
		keywordRuleMap.RuleName = UpdateKeywordRuleParams.RuleName
		keywordRuleMap.Content = val
		// 创建关键词规则
		err = db.CreateKeywordRuleMap(keywordRuleMap)
		if err != nil {
			JSONFail(ctx, OperationDBError, fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
			log.Println(fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
			return
		}
	}

	JSONSuccess(ctx, SuccessMessage)
}

// 查询关键词规则
func (controller Controller) ShowKeywordRule(ctx *gin.Context) {

	// 查询关键词规则
	keywordRule, err := db.SelectKeywordRule()
	if err != nil {
		JSONFail(ctx, OperationDBError, fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
		log.Println(fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
		return
	}

	result := make([]model.ShowKeywordRule, 0)
	for _, val := range keywordRule {
		showKeywordRule := model.ShowKeywordRule{}
		showKeywordRule.ID = val.ID
		showKeywordRule.RuleName = val.RuleName
		showKeywordRule.Switch = val.Switch
		showKeywordRule.Content = val.Content
		// 查询关键词
		keywordRuleMap, err := db.SelectKeywordRuleMapByRuleName(val.RuleName)
		if err != nil {
			JSONFail(ctx, OperationDBError, fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
			log.Println(fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
			return
		}
		showKeywordRule.Keyword = make([]string, 0)
		for _, tmp := range keywordRuleMap {
			showKeywordRule.Keyword = append(showKeywordRule.Keyword, tmp.Content)
		}
		result = append(result, showKeywordRule)
	}

	JSONSuccess(ctx, result)
}

// 删除关键词规则
func (controller Controller) DeleteKeywordRule(ctx *gin.Context) {

	// 校验请求参数
	var DeleteKeywordRuleParams model.DeleteKeywordRuleParams
	if err := ctx.ShouldBindBodyWith(&DeleteKeywordRuleParams, binding.JSON); err != nil {
		JSONFail(ctx, IllegalRequestParams, fmt.Sprintf(`%s: %s`, RequestParamsErrMessage, err.Error()))
		log.Println(fmt.Sprintf(`%s: %s`, RequestParamsErrMessage, err.Error()))
		return
	}

	// 查询关键词规则
	keywordRule, err := db.SelectKeywordRuleByRuleName(DeleteKeywordRuleParams.RuleName)
	if err != nil {
		JSONFail(ctx, OperationDBError, fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
		log.Println(fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
		return
	}

	// 删除关键词规则
	err = db.DeleteKeywordRule(keywordRule)
	if err != nil {
		JSONFail(ctx, OperationDBError, fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
		log.Println(fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
		return
	}

	// 删除规则与关键词映射
	err = db.DeleteKeywordRuleMapByRuleName(DeleteKeywordRuleParams.RuleName)
	if err != nil {
		JSONFail(ctx, OperationDBError, fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
		log.Println(fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
		return
	}

	JSONSuccess(ctx, SuccessMessage)
}

// 编辑首次回复
func (controller Controller) UpdateFirstReply(ctx *gin.Context) {

	// 校验请求参数
	var UpdateFirstReplyParams model.UpdateFirstReplyParams
	if err := ctx.ShouldBindBodyWith(&UpdateFirstReplyParams, binding.JSON); err != nil || (UpdateFirstReplyParams.OptionSwitch != 0 && UpdateFirstReplyParams.OptionSwitch != 1) {
		JSONFail(ctx, IllegalRequestParams, fmt.Sprintf(`%s: %s`, RequestParamsErrMessage, err.Error()))
		log.Println(fmt.Sprintf(`%s: %s`, RequestParamsErrMessage, err.Error()))
		return
	}

	// 查询首次回复设置
	firstReply, err := db.SelectFirstReply()
	if err != nil {
		JSONFail(ctx, OperationDBError, fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
		log.Println(fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
		return
	}

	// 更新首次回复消息
	firstReply.Message = UpdateFirstReplyParams.Message
	firstReply.OptionSwitch = UpdateFirstReplyParams.OptionSwitch
	err = db.SaveFirstReply(firstReply)
	if err != nil {
		JSONFail(ctx, OperationDBError, fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
		log.Println(fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
		return
	}

	// 删除首次回复选项消息
	if len(UpdateFirstReplyParams.OptionInfo) > 0 {
		err = db.DeleteFirstReplyOptionMessage()
		if err != nil {
			JSONFail(ctx, OperationDBError, fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
			log.Println(fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
			return
		}
	}
	// 创建首次回复选项消息
	for _, val := range UpdateFirstReplyParams.OptionInfo {
		firstReplyOptionMessage := new(model.FirstReplyOptionMessage)
		firstReplyOptionMessage.Option = val.Option
		firstReplyOptionMessage.Content = val.Content
		err = db.CreateFirstReplyOptionMessage(firstReplyOptionMessage)
		if err != nil {
			JSONFail(ctx, OperationDBError, fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
			log.Println(fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
			return
		}
	}

	JSONSuccess(ctx, SuccessMessage)
}

//// 查询首次回复
//func (controller Controller) ShowFirstReply(ctx *gin.Context) {
//
//	// 查询首次回复设置
//	firstReply, err := db.SelectFirstReply()
//	if err != nil {
//		JSONFail(ctx, OperationDBError, fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
//		log.Println(fmt.Sprintf(`%s: %s`, OperationDBErrMessage, err.Error()))
//		return
//	}
//
//
//}
