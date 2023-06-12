package postmessage

type Request struct {
	Message    string     `json:"message"`
	Attachment Attachment `json:"attachment"`
}

type Attachment struct {
	AttachmentType string  `json:"type"`
	Payload        Payload `json:"payload"`
}

type Payload struct {
	Src string `json:"src"`
}

type Response struct {
	RecipientID string `json:"recipient_id"`
	MessageID   string `json:"message_id"`
	Timestamp   int64  `json:"timestamp"`
}
