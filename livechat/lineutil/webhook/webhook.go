package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/Nhuengzii/botio-livechat-backend/livechat"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/lineutil/messagefmt"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func ValidateSignature(channelSecret string, signature string, body string) (_ bool, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("lineutil.ValidateSignature: %w", err)
		}
	}()
	decoded, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, err
	}
	hash := hmac.New(sha256.New, []byte(channelSecret))
	_, err = hash.Write([]byte(body))
	if err != nil {
		return false, err
	}
	if !hmac.Equal(decoded, hash.Sum(nil)) {
		return false, nil
	}
	return true, nil
}

type Body struct {
	Destination string           `json:"destination"` // bot user id that should receive the webhook
	Events      []*linebot.Event `json:"events"`
}

func (w *Body) HandleWebhookBodyAndExtractMessages() []*livechat.StdMessage {
	botUserID := w.Destination
	var messages []*livechat.StdMessage
	for _, event := range w.Events {
		switch event.Type {
		case linebot.EventTypeMessage:
			messages = append(messages, messagefmt.NewStdMessage(event, botUserID))
		default:
			// TODO implement user join/leave events -> updateConversationParticipants
			// info to be updated: group pic, group name, group members, and each member's name and profile pic
		}
	}
	return messages
}
