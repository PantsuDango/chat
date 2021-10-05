package model

type ChatMessageInfo struct {
	Message     string       `json:"Message"`
	MessageType string       `json:"MessageType"`
	CreateTime  string       `json:"CreateTime"`
	OptionInfo  []OptionInfo `json:"OptionInfo"`
}

type OptionInfo struct {
	Option  string `json:"Option"`
	Content string `json:"Content"`
}

type ChatMessageIPListInfo struct {
	IP          string       `json:"IP"`
	Message     string       `json:"Message"`
	MessageType string       `json:"MessageType"`
	CreateTime  string       `json:"CreateTime"`
	OptionInfo  []OptionInfo `json:"OptionInfo"`
}
