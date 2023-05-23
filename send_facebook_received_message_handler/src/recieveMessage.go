package main

type RecieveMessage struct {
	Object string         `json:"object"`
	Entry  []Notification `json:"entry"`
}

type Notification struct {
	PageID       string        `json:"id"`
	Time         int64         `json:"time"`
	MessageDatas []MessageData `json:"messaging"`
}

type MessageData struct {
	Sender    User    `json:"sender"`
	Recipient User    `json:"recipient"`
	Timestamp int64   `json:"timestamp"`
	Message   Message `json:"message"`
}

type Message struct {
	MessageID   string       `json:"mid"`
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
	ReplyTo     ReplyMessage `json:"reply_to"`
}

type User struct {
	ID string `json:"id"`
}

type Attachment struct {
	AttachmentType string      `json:"type"`
	Payload        PayloadType `json:"payload"`
}

type PayloadType struct {
	Src string `json:"src"`
}

type ReplyMessage struct {
	MessageId string `json:"messageID"`
}

type WebsocketMessage struct {
	Action  string `json:"action"`
	Message string `json:"message"`
}

type StandardMessage struct {
	ShopID         string       `json:"shopID"`
	Platform       string       `json:"platform"`
	PageID         string       `json:"pageID"`
	ConversationID string       `json:"conversationID"`
	MessageID      string       `json:"messageID"`
	Timestamp      int64        `json:"timestamp"`
	Source         Source       `json:"source"`
	Message        string       `json:"message"`
	Attachments    []Attachment `json:"attachments"`
	ReplyTo        ReplyMessage `json:"replyTo"`
}

type Source struct {
	UserID   string `json:"userID"`
	UserType string `json:"type"`
}
