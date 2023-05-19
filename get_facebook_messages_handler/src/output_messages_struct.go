package main

type OutputMessage struct {
	Messages []Message `json:"messages"`
}

type Message struct {
	MessageID      string       `json:"messageID"`
	Timestamp      int64        `json:"timestamp"`
	Source         Source       `json:"source"`
	Message        string       `json:"message"`
	AttachMents    []Attachment `json:"attachments"`
	ReplyTo        ReplyMessage `json:"replyTo"`
	ReadStatus     bool         `json:"readStatus"`
	DeliveryStatus bool         `json:"deliveryStatus"`
	UnsendStatus   bool         `json:"unsendStatus"`
}
type Attachment struct {
	AttachmentType string  `json:"type"`
	Payload        Payload `json:"payload"`
}

type ReplyMessage struct {
	MessageId string `json:"messageID"`
}

type Source struct {
	UserID   string `json:"userID"`
	UserType string `json:"type"`
}

type Payload struct {
	Src string `json:"src"`
}
