package livechat

type StdConversation struct {
	ShopID          string         `json:"shopID" bson:"shopID"`
	PageID          string         `json:"pageID" bson:"pageID"`
	ConversationID  string         `json:"conversationID" bson:"conversationID"`
	ConversationPic Payload        `json:"conversationPic" bson:"conversationPic"`
	UpdatedTime     int64          `json:"updatedTime" bson:"updatedTime"`
	Participants    []*Participant `json:"participants" bson:"participants"`
	LastActivity    string         `json:"lastActivity" bson:"lastActivity"`
	IsRead          bool           `json:"isRead" bson:"isRead"`
}

type Participant struct {
	UserID     string  `json:"userID" bson:"userID"`
	Username   string  `json:"username" bson:"username"`
	ProfilePic Payload `json:"profilePic" bson:"profilePic"`
}

// Payload here is the same as Payload in stdmessage.go
//type Payload struct {
//	Src string `json:"src" bson:"src"`
//}
