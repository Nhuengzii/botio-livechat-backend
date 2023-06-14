package postmessage

import "github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"

type Request struct {
	Message    string     `json:"message"`
	Attachment Attachment `json:"attachment"`
}

type Attachment struct {
	AttachmentType stdmessage.AttachmentType `json:"type"`
	Payload        Payload                   `json:"payload"`
}

type Payload struct {
	Src string `json:"src"`
}

type Response struct {
	RecipientID string `json:"recipient_id"`
	MessageID   string `json:"message_id"`
	Timestamp   int64  `json:"timestamp"`
}
