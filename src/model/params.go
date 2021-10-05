package model

type ShowChatMessageParams struct {
	UserType string `json:"UserType" binding:"required"`
}
