package main

import (
	"fmt"
	"github.com/Nhuengzii/botio-livechat-backend/livechat"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"reflect"
)

func newStdMessage(event *linebot.Event, botUserID string) *livechat.StdMessage {
	platform := livechat.PlatformLine
	pageID := botUserID
	shopID := "1" // TODO get from some db with botUserID?
	source := ToStdMessageSource(event.Source)
	conversationID := botUserID + ":" + source.UserID
	timestamp := event.Timestamp.UnixMilli()

	// message-type-specific fields
	var messageID string
	var message string
	var attachments []*livechat.Attachment
	var replyTo *livechat.RepliedMessage

	switch msg := event.Message.(type) {
	case *linebot.TextMessage:
		messageID = msg.ID
		message = msg.Text
		if HasLineEmojis(msg) {
			attachments = ToLineEmojiAttachments(msg) // currently nil
		}
	case *linebot.ImageMessage:
		messageID = msg.ID
		attachments = append(attachments, ToImageAttachment(msg))
	case *linebot.VideoMessage:
		messageID = msg.ID
		attachments = append(attachments, ToVideoAttachment(msg))
	case *linebot.AudioMessage:
		messageID = msg.ID
		attachments = append(attachments, ToAudioAttachment(msg))
	case *linebot.StickerMessage:
		messageID = msg.ID
		attachments = append(attachments, ToStickerAttachment(msg))
	case *linebot.LocationMessage:
		messageID = msg.ID
		message = ToLocationString(msg)
	}
	return &livechat.StdMessage{
		ShopID:         shopID,
		Platform:       platform,
		PageID:         pageID,
		ConversationID: conversationID,
		MessageID:      messageID,
		Timestamp:      timestamp,
		Source:         *source,
		Message:        message,
		Attachments:    attachments, // always nil for pure texts and locations, currently nil for texts with line emoji(s) and pure line emojis
		ReplyTo:        replyTo,     // always nil
	}
}

func ToStdMessageSource(s *linebot.EventSource) *livechat.Source {
	var userID string
	var userType livechat.UserType
	switch s.Type {
	case linebot.EventSourceTypeUser:
		userID = s.UserID
		userType = livechat.UserTypeUser
	case linebot.EventSourceTypeGroup:
		userID = s.GroupID
		userType = livechat.UserTypeGroup
	}
	return &livechat.Source{
		UserID:   userID,
		UserType: userType,
	}
}

func ToImageAttachment(m *linebot.ImageMessage) *livechat.Attachment {
	// TODO get image file from m.ID and save it to some db
	return &livechat.Attachment{
		AttachmentType: livechat.AttachmentTypeImage,
		Payload: livechat.Payload{
			Src: "", // TODO get url of the image stored in some db
		}}
}

func ToVideoAttachment(m *linebot.VideoMessage) *livechat.Attachment {
	// TODO get video file from m.ID and save it to some db
	return &livechat.Attachment{
		AttachmentType: livechat.AttachmentTypeVideo,
		Payload: livechat.Payload{
			Src: "", // TODO get url of the video stored in some db
		}}
}

func ToAudioAttachment(m *linebot.AudioMessage) *livechat.Attachment {
	// TODO get audio file from m.ID and save it to some db
	return &livechat.Attachment{
		AttachmentType: livechat.AttachmentTypeAudio,
		Payload: livechat.Payload{
			Src: "", // TODO get url of the audio stored in some db
		}}
}

func ToStickerAttachment(m *linebot.StickerMessage) *livechat.Attachment {
	return &livechat.Attachment{
		AttachmentType: livechat.AttachmentTypeSticker,
		Payload: livechat.Payload{
			Src: ToStickerURL(m),
		}}
}

func ToStickerURL(m *linebot.StickerMessage) string {
	return fmt.Sprintf("https://stickershop.line-scdn.net/stickershop/v1/sticker/%s/android/sticker.png", m.StickerID)
}

func HasLineEmojis(m *linebot.TextMessage) bool {
	v := reflect.ValueOf(m).Elem().FieldByName("Emojis")
	return v != reflect.Value{}
}

func ToLineEmojiAttachments(m *linebot.TextMessage) []*livechat.Attachment {
	var attachments []*livechat.Attachment
	// TODO implement me
	return attachments
}

func ToLineEmojiURL(e *linebot.Emoji) string {
	return fmt.Sprintf("https://stickershop.line-scdn.net/sticonshop/v1/sticon/%s/android/%s.png", e.ProductID, e.EmojiID)
}

func ToLocationString(m *linebot.LocationMessage) string {
	return fmt.Sprintf("Title: %s\nAddress: %s\nLatitude: %f\nLongitude: %f", m.Title, m.Address, m.Latitude, m.Longitude)
}
