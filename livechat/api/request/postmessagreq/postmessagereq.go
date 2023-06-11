package postmessagreq

type Request struct {
	Message    string     `json:"message"`
	Attachment Attachment `json:"attachment"`
}

type Attachment struct {
	AttachmentType string      `json:"type"`
	Payload        PayloadType `json:"payload"`
}

type PayloadType struct {
	Src string `json:"src"`
}
