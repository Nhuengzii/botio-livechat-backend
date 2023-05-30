package main

type botioConversation struct {
	ShopID          string        `bson:"shopID"`
	PageID          string        `bson:"pageID"`
	ConversationID  string        `bson:"conversationID"`
	ConversationPic payload       `bson:"conversationPic"`
	UpdatedTime     int64         `bson:"updatedTime"`
	Participants    []participant `bson:"participants"`
	LastActivity    string        `bson:"lastActivity"`
	IsRead          bool          `bson:"isRead"` // this field is always false for LINE (unsupported)
}

type participant struct {
	UserID     string  `bson:"userID"`
	Username   string  `bson:"username"`
	ProfilePic payload `bson:"profilePic"`
}

type payload struct {
	Src string `json:"src"`
}
