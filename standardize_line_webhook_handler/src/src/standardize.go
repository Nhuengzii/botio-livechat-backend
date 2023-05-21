package main

import (
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

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
	UserType string `json:"userType"`
}

type Attachment struct {
	AttachmentType string  `json:"attachmentType"`
	Payload        Payload `json:"payload"`
}

type Payload struct {
	Src string `json:"url"`
}

type ReplyMessage struct {
	MessageID string `json:"messageID"`
}

type WebhookBody struct {
	Destination string           `json:"destination"`
	Events      []*linebot.Event `json:"events"`
}

func (wb WebhookBody) toStandardMessages() []StandardMessage {
	botUserID := wb.Destination
	var standardMessages []StandardMessage
	for _, event := range wb.Events {
		eventSource := extractSource(event.Source)
		if event.Type != linebot.EventTypeMessage {
			continue
		}
		standardMessage := StandardMessage{
			ShopID:         "0", // get from redis key: botUserID
			Platform:       "LINE",
			PageID:         botUserID,
			ConversationID: botUserID + ":" + eventSource.UserID,
			//MessageID:
			Timestamp: event.Timestamp.UnixMilli(),
			Source:    eventSource,
			//Message:
			Attachments: []Attachment{},
			//ReplyTo:
		}
		switch message := event.Message.(type) {
		case *linebot.TextMessage:
			standardMessage.MessageID = message.ID
			standardMessage.Message = message.Text
		case *linebot.ImageMessage:
			standardMessage.MessageID = message.ID
		case *linebot.VideoMessage:
			standardMessage.MessageID = message.ID
		case *linebot.AudioMessage:
			standardMessage.MessageID = message.ID
		case *linebot.FileMessage:
			standardMessage.MessageID = message.ID
		case *linebot.LocationMessage:
			standardMessage.MessageID = message.ID
		case *linebot.StickerMessage:
			standardMessage.MessageID = message.ID
			standardMessage.Attachments = append(standardMessage.Attachments, Attachment{
				AttachmentType: "sticker",
				Payload: Payload{
					Src: "https://stickershop.line-scdn.net/stickershop/v1/sticker/" + message.StickerID + "/android/sticker.png",
				},
			})
		}
		standardMessages = append(standardMessages, standardMessage)
	}
	return standardMessages
}

func extractSource(eventSource *linebot.EventSource) Source {
	source := Source{}
	switch eventSource.Type {
	case linebot.EventSourceTypeUser:
		source.UserID = eventSource.UserID
		source.UserType = "user"
	case linebot.EventSourceTypeGroup:
		source.UserID = eventSource.GroupID
		source.UserType = "group"
	case linebot.EventSourceTypeRoom:
		source.UserID = eventSource.RoomID
		source.UserType = "room"
	}
	return source
}
