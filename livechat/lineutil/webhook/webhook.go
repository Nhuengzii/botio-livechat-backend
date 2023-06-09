package webhook

import (
	"github.com/Nhuengzii/botio-livechat-backend/livechat"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/lineutil/messagefmt"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type Body struct {
	Destination string           `json:"destination"` // bot user id that should receive the webhook
	Events      []*linebot.Event `json:"events"`
}

func (w *Body) HandleWebhookBody() []*livechat.StdMessage {
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
