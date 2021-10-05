package model

type ShowChatMessageParams struct {
	UserType string `json:"UserType" binding:"required"`
}

type UpdateIpContentMapParams struct {
	IP        string `json:"IP"        binding:"required"`
	IpContent string `json:"IpContent" binding:"required"`
}
