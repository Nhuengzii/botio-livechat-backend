package main

type facebook_message struct {
	ShopId         string       `bson:"shopID"`
	Platform       string       `bson:"platform"`
	PageId         string       `bson:"pageId"`
	ConversationId string       `bson:"conversationID"`
	MessageId      string       `bson:"messageID"`
	TimeStamp      int64        `bson:"timestamp"`
	Source         Source       `bson:"source"`
	Message        string       `bson:"message"`
	Attachments    []Attachment `bson:"attachments"`
	ReplyTo        Reply        `bson:"replyTo"`
	ReadStatus     bool         `bson:"readStatus"`
	DeliveryStatus bool         `bson:"deliveryStatus"`
	UnsendStatus   bool         `bson:"unsendStatus"`
}

type Source struct {
	UserId   string `bson:"userID"`
	UserType string `bson:"userType"`
}

type Attachment struct {
	AttachmentType string `bson:"attachmentType"`
	Payload        struct {
		Src string `bson:"src"`
	} `bson:"payload"`
}

type Reply struct {
	MessageId string `bson:"messageID"`
}
