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
	IPContent   string       `json:"IPContent"`
	Message     string       `json:"Message"`
	MessageType string       `json:"MessageType"`
	CreateTime  string       `json:"CreateTime"`
	OptionInfo  []OptionInfo `json:"OptionInfo"`
}

type ShowKeywordRule struct {
	ID       int      `json:"ID"`
	RuleName string   `json:"RuleName"`
	Switch   int      `json:"Switch"`
	Content  string   `json:"Content"`
	Keyword  []string `json:"Keyword"`
}

type ShowFirstReply struct {
	Message      string       `json:"Message"`
	OptionSwitch int          `json:"OptionSwitch"`
	OptionInfo   []OptionInfo `json:"OptionInfo"`
}
