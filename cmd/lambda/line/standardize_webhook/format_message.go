package main

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func (c *config) newStdMessage(shopID string, pageID string, event *linebot.Event) (*stdmessage.StdMessage, error) {
	platform := stdmessage.PlatformLine
	source, err := toStdMessageSource(event.Source)
	if err != nil {
		return nil, fmt.Errorf("newStdMessage: %w", err)
	}
	conversationID := source.UserID
	timestamp := event.Timestamp.UnixMilli()

	// message-type-specific fields
	var messageID string
	var message string
	var attachments []stdmessage.Attachment
	var replyTo *stdmessage.RepliedMessage

	switch msg := event.Message.(type) {
	case *linebot.TextMessage:
		messageID = msg.ID
		message = msg.Text
		if hasLineEmojis(msg) {
			attachments = toLineEmojiAttachments(msg) // currently nil
		}
	case *linebot.ImageMessage:
		messageID = msg.ID
		attachments = append(attachments, toImageAttachment(msg))
	case *linebot.VideoMessage:
		messageID = msg.ID
		attachments = append(attachments, toVideoAttachment(msg))
	case *linebot.AudioMessage:
		messageID = msg.ID
		attachments = append(attachments, toAudioAttachment(msg))
	case *linebot.StickerMessage:
		messageID = msg.ID
		attachments = append(attachments, toStickerAttachment(msg))
	case *linebot.LocationMessage:
		messageID = msg.ID
		message = toLocationString(msg)
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
	}, nil
}

var errMessageSourceUnsupported = errors.New("message source unsupported")

func toStdMessageSource(s *linebot.EventSource) (*stdmessage.Source, error) {
	if s.Type != linebot.EventSourceTypeUser {
		return nil, fmt.Errorf("toStdMessageSource: %w", errMessageSourceUnsupported)
	}
	return &stdmessage.Source{
		UserID:   s.UserID,
		UserType: stdmessage.UserTypeUser,
	}, nil
}

func toImageAttachment(m *linebot.ImageMessage) stdmessage.Attachment {
	// TODO get image file from m.ID and save it to some db
	return stdmessage.Attachment{
		AttachmentType: stdmessage.AttachmentTypeImage,
		Payload: stdmessage.Payload{
			Src: "", // TODO get url of the image stored in some db
		},
	}
}

func toVideoAttachment(m *linebot.VideoMessage) stdmessage.Attachment {
	// TODO get video file from m.ID and save it to some db
	return stdmessage.Attachment{
		AttachmentType: stdmessage.AttachmentTypeVideo,
		Payload: stdmessage.Payload{
			Src: "", // TODO get url of the video stored in some db
		},
	}
}

func toAudioAttachment(m *linebot.AudioMessage) stdmessage.Attachment {
	// TODO get audio file from m.ID and save it to some db
	return stdmessage.Attachment{
		AttachmentType: stdmessage.AttachmentTypeAudio,
		Payload: stdmessage.Payload{
			Src: "", // TODO get url of the audio stored in some db
		},
	}
}

func toStickerAttachment(m *linebot.StickerMessage) stdmessage.Attachment {
	return stdmessage.Attachment{
		AttachmentType: stdmessage.AttachmentTypeSticker,
		Payload: stdmessage.Payload{
			Src: toStickerURL(m),
		},
	}
}

func toStickerURL(m *linebot.StickerMessage) string {
	return fmt.Sprintf("https://stickershop.line-scdn.net/stickershop/v1/sticker/%s/android/sticker.png", m.StickerID)
}

func hasLineEmojis(m *linebot.TextMessage) bool {
	v := reflect.ValueOf(m).Elem().FieldByName("Emojis")
	return v != reflect.Value{}
}

func toLineEmojiAttachments(m *linebot.TextMessage) []stdmessage.Attachment {
	var attachments []stdmessage.Attachment
	// TODO implement me
	return attachments
}

func toLineEmojiURL(e *linebot.Emoji) string {
	return fmt.Sprintf("https://stickershop.line-scdn.net/sticonshop/v1/sticon/%s/android/%s.png", e.ProductID, e.EmojiID)
}

func toLocationString(m *linebot.LocationMessage) string {
	return fmt.Sprintf("Title: %s\nAddress: %s\nLatitude: %f\nLongitude: %f", m.Title, m.Address, m.Latitude, m.Longitude)
}
