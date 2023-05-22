package main

type RequestMessage struct {
	Message     string       `json:"message"`
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	AttachmentType string      `json:"type"`
	Payload        PayloadType `json:"payload"`
}

type PayloadType struct {
	Src string `json:"src"`
}
