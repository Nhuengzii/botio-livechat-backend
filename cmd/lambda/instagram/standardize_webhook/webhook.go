package main

type ReceiveWebhook struct {
	Object  string  `json:"object"`
	Entries []Entry `json:"entry"`
}

type Entry struct {
	PageID     string      `json:"id"`
	Time       int64       `json:"time"`
	Messagings []Messaging `json:"messaging"`
}

type Messaging struct {
	Sender    User    `json:"sender"`
	Recipient User    `json:"recipient"`
	Timestamp int64   `json:"timestamp"`
	Message   Message `json:"message"`
}

type Message struct {
	IsEcho        bool         `json:"is_echo"`
	IsDeleted     bool         `json:"is_deleted"`
	IsUnsupported bool         `json:"is_unsupported"`
	MessageID     string       `json:"mid"`
	Text          string       `json:"text"`
	Attachments   []Attachment `json:"attachments"`
	ReplyTo       ReplyMessage `json:"reply_to"`
}

type User struct {
	ID string `json:"id"`
}

type Attachment struct {
	AttachmentType string `json:"type"`
	Payload        any    `json:"payload"`
}

type BasicPayload struct {
	Src string `json:"url"`
}
type TemplatePayload struct {
	Generic map[string]any `json:"generic"`
	Product map[string]any `json:"product"`
}

type ReplyMessage struct {
	MessageId string `json:"mid"`
}
