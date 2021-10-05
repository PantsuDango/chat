package db

import (
	model2 "chat/src/model"
)

// 查询聊天消息记录
func SelectChatMessage(ip string) (chatMessage []*model2.ChatMessage, err error) {

	chatMessage = make([]*model2.ChatMessage, 0)
	err = exeDB.Where(map[string]interface{}{"ip": ip}).Find(&chatMessage).Error

	return
}

// 查询聊天消息记录
func SelectAllChatMessage() (chatMessage []*model2.ChatMessage, err error) {

	chatMessage = make([]*model2.ChatMessage, 0)
	err = exeDB.Find(&chatMessage).Error

	return
}

// 创建聊天消息记录
func CreateChatMessage(chatMessage *model2.ChatMessage) (err error) {

	err = exeDB.Create(&chatMessage).Error

	return
}

// 查询首次回复设置
func SelectFirstReply() (firstReply *model2.FirstReply, err error) {

	firstReply = new(model2.FirstReply)
	err = exeDB.First(&firstReply).Error

	return
}

// 查询首次回复选项信息
func SelectFirstReplyOptionMessage() (firstReplyOptionMessage []*model2.FirstReplyOptionMessage, err error) {

	firstReplyOptionMessage = make([]*model2.FirstReplyOptionMessage, 0)
	err = exeDB.First(&firstReplyOptionMessage).Error

	return
}
