package main

import (
	"errors"
	"fmt"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/storage/amazons3"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func (c *config) newStdMessage(shopID string, pageID string, bot *linebot.Client, uploader amazons3.Uploader, event *linebot.Event) (_ *stdmessage.StdMessage, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("newStdMessage: %w", err)
		}
	}()

	platform := stdmessage.PlatformLine
	if event.Source.Type != linebot.EventSourceTypeUser {
		return nil, errors.New("event source type unsupported")
	}
	source := stdmessage.Source{
		UserID:   event.Source.UserID,
		UserType: stdmessage.UserTypeUser,
	}
	conversationID := source.UserID
	timestamp := event.Timestamp.UnixMilli()

	// message-type-specific fields
	var messageID string
	var message string
	attachments := []stdmessage.Attachment{}
	var replyTo *stdmessage.RepliedMessage

	switch msg := event.Message.(type) {
	case *linebot.TextMessage:
		messageID = msg.ID
		message = msg.Text
		if hasLineEmojis(msg) {
			attachments = toLineEmojiAttachments(msg) // currently empty slice
		}
	case *linebot.ImageMessage:
		messageID = msg.ID
		location, err := getAndUploadMessageContent(bot, uploader, messageID)
		if err != nil {
			return nil, err
		}
		attachments = append(attachments, stdmessage.Attachment{
			AttachmentType: stdmessage.AttachmentTypeImage,
			Payload: stdmessage.Payload{
				Src: location,
			},
		})
	case *linebot.VideoMessage:
		messageID = msg.ID
		location, err := getAndUploadMessageContent(bot, uploader, messageID)
		if err != nil {
			return nil, err
		}
		attachments = append(attachments, stdmessage.Attachment{
			AttachmentType: stdmessage.AttachmentTypeVideo,
			Payload: stdmessage.Payload{
				Src: location,
			},
		})
	case *linebot.AudioMessage:
		messageID = msg.ID
		location, err := getAndUploadMessageContent(bot, uploader, messageID)
		if err != nil {
			return nil, err
		}
		attachments = append(attachments, stdmessage.Attachment{
			AttachmentType: stdmessage.AttachmentTypeAudio,
			Payload: stdmessage.Payload{
				Src: location,
			},
		})
	case *linebot.StickerMessage:
		messageID = msg.ID
		attachments = append(attachments, stdmessage.Attachment{
			AttachmentType: stdmessage.AttachmentTypeSticker,
			Payload: stdmessage.Payload{
				Src: toStickerURL(msg),
			},
		})
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
		Source:         source,
		Message:        message,
		Attachments:    attachments, // always empty for pure texts and locations, currently empty for texts with line emoji(s) and pure line emojis
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

func getAndUploadMessageContent(bot *linebot.Client, uploader amazons3.Uploader, messageID string) (_ string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("getAndUploadMessageContent: %w", err)
		}
	}()
	response, err := bot.GetMessageContent(messageID).Do()
	if err != nil {
		return "", err
	}
	file := response.Content
	location, err := uploader.UploadFile(file)
	if err != nil {
		return "", err
	}
	return location, nil
}

func toStickerURL(m *linebot.StickerMessage) string {
	return fmt.Sprintf("https://stickershop.line-scdn.net/stickershop/v1/sticker/%s/android/sticker.png", m.StickerID)
}

func hasLineEmojis(m *linebot.TextMessage) bool {
	return m.Emojis != nil
}

func toLineEmojiAttachments(m *linebot.TextMessage) []stdmessage.Attachment {
	attachments := []stdmessage.Attachment{}
	// TODO implement me
	return attachments
}

func toLineEmojiURL(e *linebot.Emoji) string {
	return fmt.Sprintf("https://stickershop.line-scdn.net/sticonshop/v1/sticon/%s/android/%s.png", e.ProductID, e.EmojiID)
}

func toLocationString(m *linebot.LocationMessage) string {
	return fmt.Sprintf("Title: %s\nAddress: %s\nLatitude: %f\nLongitude: %f", m.Title, m.Address, m.Latitude, m.Longitude)
}
