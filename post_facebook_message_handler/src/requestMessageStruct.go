package main

type RequestMessage struct {
	Message    string     `json:"message"`
	Attachment Attachment `json:"attachment"`
}

type Attachment struct {
	AttachmentType string      `json:"type" bson:"type"`
	Payload        PayloadType `json:"payload" bson:"payload"`
}

type PayloadType struct {
	Src string `json:"src" bson:"src"`
}
