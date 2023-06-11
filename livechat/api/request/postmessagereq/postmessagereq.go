package postmessagereq

type Request struct {
	Message    string     `json:"message"`
	Attachment Attachment `json:"attachment"` // Why is this not Attachments []Attachment `json:"attachments"`?
}

type Attachment struct {
	AttachmentType string  `json:"type"`
	Payload        Payload `json:"payload"`
}

type Payload struct {
	Src string `json:"src"`
}
