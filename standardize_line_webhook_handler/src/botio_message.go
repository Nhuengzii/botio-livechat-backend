package main

import "github.com/line/line-bot-sdk-go/v7/linebot"

type botioMessage struct {
	ShopID         string        `json:"shopID"`
	Platform       platform      `json:"platform"`
	PageID         string        `json:"pageID"`
	ConversationID string        `json:"conversationID"`
	MessageID      string        `json:"messageID"`
	Timestamp      int64         `json:"timestamp"`
	Source         source        `json:"source"`
	Message        string        `json:"message"`
	Attachments    []attachment  `json:"attachments"`
	ReplyTo        *replyMessage `json:"replyTo"`
}

type platform string

const (
	platformLine platform = "line"
)

type source struct {
	UserID   string   `json:"userID"`
	UserType userType `json:"userType"`
}

type userType string

const (
	userTypeUser  userType = "user"
	userTypeGroup userType = "group"
	userTypeRoom  userType = "room"
)

// the currne attachment structure doesn't allow
// more than one line-emojis on their own or
// line-emojis embedding in text messages
type attachment struct {
	AttachmentType attachmentType `json:"attachmentType"`
	Payload        payload        `json:"payload"`
}

type attachmentType string

const (
	// attachmentTypeLineEmoji attachmentType = "lineEmoji"
	// attachmentTypeImage     attachmentType = "image"
	// attachmentTypeVideo     attachmentType = "video"
	// attachmentTypeAudio     attachmentType = "audio"
	attachmentTypeSticker attachmentType = "sticker"
)

type payload struct {
	Src string `json:"src"`
}

type replyMessage struct {
	MessageID string `json:"messageID"`
}

func newBotioMessageSource(es *linebot.EventSource) source {
	var uID string
	var uType userType
	switch es.Type {
	case linebot.EventSourceTypeUser:
		uID = es.UserID
		uType = userTypeUser
	case linebot.EventSourceTypeGroup:
		uID = es.GroupID
		uType = userTypeGroup
	case linebot.EventSourceTypeRoom:
		uID = es.RoomID
		uType = userTypeRoom
	}
	return source{
		UserID:   uID,
		UserType: uType,
	}
}
