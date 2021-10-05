package db

import (
	"chat/src/model"
	"time"
)

const (
	TimeFormat = "2006-01-02 15:04:05"
)

// 查询聊天消息记录
func SelectChatMessage(ip string) (chatMessage []*model.ChatMessage, err error) {

	now := time.Now()
	dd, _ := time.ParseDuration("-168h")
	end := now.Format(TimeFormat)
	start := now.Add(dd)

	chatMessage = make([]*model.ChatMessage, 0)
	err = exeDB.Where(`ip = ? and createtime between ? and ?`, ip, start, end).Find(&chatMessage).Error

	return
}

// 查询聊天消息记录
func SelectAllChatMessage() (chatMessage []*model.ChatMessage, err error) {

	now := time.Now()
	dd, _ := time.ParseDuration("-168h")
	end := now.Format(TimeFormat)
	start := now.Add(dd)

	chatMessage = make([]*model.ChatMessage, 0)
	err = exeDB.Where(`createtime between ? and ?`, start, end).Find(&chatMessage).Error

	return
}

// 查询聊天消息记录
func SelectAllChatMessageIP() (chatMessage []*model.ChatMessage, err error) {

	now := time.Now()
	dd, _ := time.ParseDuration("-168h")
	end := now.Format(TimeFormat)
	start := now.Add(dd)

	chatMessage = make([]*model.ChatMessage, 0)
	err = exeDB.Where(`createtime between ? and ?`, start, end).Order(`createtime desc`).Find(&chatMessage).Error

	return
}

// 创建聊天消息记录
func CreateChatMessage(chatMessage *model.ChatMessage) (err error) {

	err = exeDB.Create(&chatMessage).Error
	return
}

// 查询首次回复设置
func SelectFirstReply() (firstReply *model.FirstReply, err error) {

	firstReply = new(model.FirstReply)
	err = exeDB.First(&firstReply).Error
	return
}

// 查询首次回复选项信息
func SelectFirstReplyOptionMessage() (firstReplyOptionMessage []*model.FirstReplyOptionMessage, err error) {

	firstReplyOptionMessage = make([]*model.FirstReplyOptionMessage, 0)
	err = exeDB.First(&firstReplyOptionMessage).Error
	return
}

// 查询IP备注
func SelectIpContentMapByIP(ip string) (ipContentMap *model.IpContentMap, err error) {

	ipContentMap = new(model.IpContentMap)
	err = exeDB.Where(map[string]interface{}{"ip": ip}).First(&ipContentMap).Error
	return
}

// 更新IP备注
func SaveIpContentMap(ipContentMap *model.IpContentMap) (err error) {

	err = exeDB.Save(&ipContentMap).Error
	return
}

// 查询IP备注
func SelectIpContentMap() (ipContentMap []*model.IpContentMap, err error) {

	err = exeDB.Find(&ipContentMap).Error
	return
}
