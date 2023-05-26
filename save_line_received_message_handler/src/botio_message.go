package main

type botioMessage struct {
	ShopID         string        `json:"shopID" bson:"shopID"`
	Platform       platform      `json:"platform" bson:"platform"`
	PageID         string        `json:"pageID" bson:"pageID"`
	ConversationID string        `json:"conversationID" bson:"conversationID"`
	MessageID      string        `json:"messageID" bson:"messageID"`
	Timestamp      int64         `json:"timestamp" bson:"timestamp"`
	Source         source        `json:"source" bson:"source"`
	Message        string        `json:"message" bson:"message"`
	Attachments    []attachment  `json:"attachments" bson:"attachments"`
	ReplyTo        *replyMessage `json:"replyTo" bson:"replyTo"`
}

type platform string

// const (
// 	platformLine platform = "line"
// )

type source struct {
	UserID   string   `json:"userID" bson:"userID"`
	UserType userType `json:"userType" bson:"userType"`
}

type userType string

// const (
// 	userTypeUser  userType = "user"
// 	userTypeGroup userType = "group"
// 	userTypeRoom  userType = "room"
// )

// the currne attachment structure doesn't allow
// more than one line-emojis on their own or
// line-emojis embedding in text messages
type attachment struct {
	AttachmentType attachmentType `json:"attachmentType" bson:"attachmentType"`
	Payload        payload        `json:"payload" bson:"payload"`
}

type attachmentType string

// const (
// 	// attachmentTypeLineEmoji attachmentType = "lineEmoji"
// 	// attachmentTypeImage     attachmentType = "image"
// 	// attachmentTypeVideo     attachmentType = "video"
// 	// attachmentTypeAudio     attachmentType = "audio"
// 	attachmentTypeSticker attachmentType = "sticker"
// )

type payload struct {
	Src string `json:"src"`
}

type replyMessage struct {
	MessageID string `json:"messageID" bson:"messageID"`
}
