package model

type ShowChatMessageParams struct {
	UserType string `json:"UserType" binding:"required"`
	IP       string `json:"IP"`
}

type UpdateIpContentMapParams struct {
	IP        string `json:"IP"        binding:"required"`
	IpContent string `json:"IpContent" binding:"required"`
}

type SendChatMessageParams struct {
	IP          string `json:"IP"`
	Message     string `json:"Message"     binding:"required"`
	MessageType string `json:"MessageType" binding:"required"`
}
