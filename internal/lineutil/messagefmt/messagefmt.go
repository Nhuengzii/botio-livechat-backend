package messagefmt

import (
	"fmt"
	"github.com/Nhuengzii/botio-livechat-backend/pkg/stdmessage"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"reflect"
)

func NewStdMessage(event *linebot.Event, botUserID string) *stdmessage.StdMessage {
	platform := stdmessage.PlatformLine
	pageID := botUserID
	shopID := "1" // TODO get from some db with botUserID?
	source := ToStdMessageSource(event.Source)
	conversationID := botUserID + ":" + source.UserID
	timestamp := event.Timestamp.UnixMilli()

	// message-type-specific fields
	var messageID string
	var message string
	var attachments []*stdmessage.Attachment
	var replyTo *stdmessage.RepliedMessage

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
	return &stdmessage.StdMessage{
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

func ToStdMessageSource(s *linebot.EventSource) *stdmessage.Source {
	var userID string
	var userType stdmessage.UserType
	switch s.Type {
	case linebot.EventSourceTypeUser:
		userID = s.UserID
		userType = stdmessage.UserTypeUser
	case linebot.EventSourceTypeGroup:
		userID = s.GroupID
		userType = stdmessage.UserTypeGroup
	}
	return &stdmessage.Source{
		UserID:   userID,
		UserType: userType,
	}
}

func ToImageAttachment(m *linebot.ImageMessage) *stdmessage.Attachment {
	// TODO get image file from m.ID and save it to some db
	return &stdmessage.Attachment{
		AttachmentType: stdmessage.AttachmentTypeImage,
		Payload: stdmessage.Payload{
			Src: "", // TODO get url of the image stored in some db
		}}
}

func ToVideoAttachment(m *linebot.VideoMessage) *stdmessage.Attachment {
	// TODO get video file from m.ID and save it to some db
	return &stdmessage.Attachment{
		AttachmentType: stdmessage.AttachmentTypeVideo,
		Payload: stdmessage.Payload{
			Src: "", // TODO get url of the video stored in some db
		}}
}

func ToAudioAttachment(m *linebot.AudioMessage) *stdmessage.Attachment {
	// TODO get audio file from m.ID and save it to some db
	return &stdmessage.Attachment{
		AttachmentType: stdmessage.AttachmentTypeAudio,
		Payload: stdmessage.Payload{
			Src: "", // TODO get url of the audio stored in some db
		}}
}

func ToStickerAttachment(m *linebot.StickerMessage) *stdmessage.Attachment {
	return &stdmessage.Attachment{
		AttachmentType: stdmessage.AttachmentTypeSticker,
		Payload: stdmessage.Payload{
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

func ToLineEmojiAttachments(m *linebot.TextMessage) []*stdmessage.Attachment {
	var attachments []*stdmessage.Attachment
	// TODO implement me
	return attachments
}

func ToLineEmojiURL(e *linebot.Emoji) string {
	return fmt.Sprintf("https://stickershop.line-scdn.net/sticonshop/v1/sticon/%s/android/%s.png", e.ProductID, e.EmojiID)
}

func ToLocationString(m *linebot.LocationMessage) string {
	return fmt.Sprintf("Title: %s\nAddress: %s\nLatitude: %f\nLongitude: %f", m.Title, m.Address, m.Latitude, m.Longitude)
}
