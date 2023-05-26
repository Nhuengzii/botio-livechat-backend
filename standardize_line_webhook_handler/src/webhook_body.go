package main

import (
	"encoding/json"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type webhookBody struct {
	Destination string           `json:"destination"` // bot user id that should receive the webhook
	Events      []*linebot.Event `json:"events"`
}

func parseWebhookBody(body string) (*webhookBody, error) {
	wb := &webhookBody{}
	if err := json.Unmarshal([]byte(body), wb); err != nil {
		return nil, &parseWebhookBodyError{
			message: "couldn't parse webhook body",
			err:     err,
		}
	}
	return wb, nil
}

func (wb *webhookBody) toBotioMessages() []*botioMessage {
	botUserID := wb.Destination
	var botioMessages = []*botioMessage{}
	for _, event := range wb.Events {
		if event.Type != linebot.EventTypeMessage {
			continue
		}
		platform := platformLine
		pageID := botUserID
		shopID := "1" // get from pageID?
		source := newBotioMessageSource(event.Source)
		conversationID := botUserID + ":" + string(source.UserType) + ":" + source.UserID
		timestamp := event.Timestamp.UnixMilli()

		var messageID string
		var message string
		var attachments = []attachment{}
		var replyTo *replyMessage

		switch m := event.Message.(type) {
		case *linebot.TextMessage:
			messageID = m.ID
			message = m.Text
			if hasLineEmojis(m) {
				attachments = toLineEmojisBotioAttachments(m) // currently returns empty []attachment{}
			}
		case *linebot.ImageMessage:
			messageID = m.ID
		case *linebot.VideoMessage:
			messageID = m.ID
		case *linebot.AudioMessage:
			messageID = m.ID
		case *linebot.LocationMessage:
			messageID = m.ID
			message = getLocationString(m)
		case *linebot.StickerMessage:
			messageID = m.ID
			attachments = toStickerBotioAttachments(m)
		}

		botioMessages = append(botioMessages,
			&botioMessage{
				ShopID:         shopID,
				Platform:       platform,
				PageID:         pageID,
				ConversationID: conversationID,
				MessageID:      messageID,
				Timestamp:      timestamp,
				Source:         source,
				Message:        message,
				Attachments:    attachments,
				ReplyTo:        replyTo,
			})
	}

	return botioMessages
}

type parseWebhookBodyError struct {
	message string
	err     error
}

func (e *parseWebhookBodyError) Error() string {
	return e.message + ": " + e.err.Error()
}