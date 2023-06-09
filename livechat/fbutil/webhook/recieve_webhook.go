package webhook

type ReceiveWebhook struct {
	Object string         `json:"object"`
	Entry  []Notification `json:"entry"`
}

type Notification struct {
	PageID       string        `json:"id"`
	Time         int64         `json:"time"`
	MessageDatas []MessageData `json:"messaging"`
}

type MessageData struct {
	Sender    User     `json:"sender"`
	Recipient User     `json:"recipient"`
	Timestamp int64    `json:"timestamp"`
	Message   Message  `json:"message"`
	Delivery  Delivery `json:"delivery"`
}

type Delivery struct {
	MessageIDs []string `json:"mids"`
	Watermark  int64    `json:"watermark"` // all messages before watermark timestamp was sent
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
	Src string `json:"url"`
}

type ReplyMessage struct {
	MessageId string `json:"messageID"`
}
