package model

import "time"

type ChatMessage struct {
	ID          int       `json:"ID"          gorm:"column:id"`
	IP          string    `json:"IP"          gorm:"column:ip"`
	Message     string    `json:"Message"     gorm:"column:message"`
	MessageType string    `json:"MessageType" gorm:"column:message_type"`
	CreatedAt   time.Time `json:"CreateTime"  gorm:"column:createtime"`
}

func (ChatMessage) TableName() string {
	return "chat_message"
}

type FirstReply struct {
	ID           int       `json:"ID"          gorm:"column:id"`
	Message      string    `json:"Message"      gorm:"column:message"`
	OptionSwitch int       `json:"OptionSwitch" gorm:"column:option_switch"`
	CreatedAt    time.Time `json:"CreateTime"   gorm:"column:createtime"`
	UpdatedAt    time.Time `json:"UpdateTime"   gorm:"column:lastupdate"`
}

func (FirstReply) TableName() string {
	return "first_reply"
}

type FirstReplyOptionMessage struct {
	ID        int       `json:"ID"          gorm:"column:id"`
	Option    string    `json:"Option"      gorm:"column:option"`
	Content   string    `json:"Content"     gorm:"column:content"`
	CreatedAt time.Time `json:"CreateTime"  gorm:"column:createtime"`
	UpdatedAt time.Time `json:"UpdateTime"  gorm:"column:lastupdate"`
}

func (FirstReplyOptionMessage) TableName() string {
	return "first_reply_option_message"
}

type IpContentMap struct {
	ID        int       `json:"ID"          gorm:"column:id"`
	IP        string    `json:"IP"          gorm:"column:ip"`
	Content   string    `json:"Content"     gorm:"column:content"`
	CreatedAt time.Time `json:"CreateTime"  gorm:"column:createtime"`
	UpdatedAt time.Time `json:"UpdateTime"  gorm:"column:lastupdate"`
}

func (IpContentMap) TableName() string {
	return "ip_content_map"
}
