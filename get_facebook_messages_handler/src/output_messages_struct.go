package main

type OutputMessage struct {
	Messages []Message `json:"messages" bson:"messages"`
}

type Message struct {
	MessageID      string       `json:"messageID" bson:"messageID"`
	Timestamp      int64        `json:"timestamp" bson:"timestamp"`
	Source         Source       `json:"source" bson:"source"`
	Message        string       `json:"message" bson:"message"`
	AttachMents    []Attachment `json:"attachments" bson:"attachments"`
	ReplyTo        ReplyMessage `json:"replyTo" bson:"replyTo"`
	ReadStatus     bool         `json:"readStatus" bson:"readStatus"`
	DeliveryStatus bool         `json:"deliveryStatus" bson:"deliveryStatus"`
	UnsendStatus   bool         `json:"unsendStatus" bson:"unsendStatus"`
}
type Attachment struct {
	AttachmentType string  `json:"type" bson:"type"`
	Payload        Payload `json:"payload" bson:"payload"`
}

type ReplyMessage struct {
	MessageId string `json:"messageID" bson:"messageID"`
}

type Source struct {
	UserID   string `json:"userID" bson:"userID"`
	UserType string `json:"type" bson:"type"`
}

type Payload struct {
	Src string `json:"src" bson:"src"`
}
