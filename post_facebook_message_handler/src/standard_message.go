package main

type StandardMessage struct {
	ShopID         string       `bson:"shopID"`
	Platform       string       `bson:"platform"`
	PageID         string       `bson:"pageID"`
	ConversationID string       `bson:"conversationID"`
	MessageID      string       `bson:"messageID"`
	Timestamp      int64        `bson:"timestamp"`
	Source         Source       `bson:"source"`
	Message        string       `bson:"message"`
	Attachments    []Attachment `bson:"attachments"`
	ReplyTo        ReplyMessage `bson:"replyTo"`
}

type Source struct {
	UserID   string `bson:"userID"`
	UserType string `bson:"type"`
}

type ReplyMessage struct {
	MessageId string `bson:"messageID"`
}