package postmessage

import "github.com/line/line-bot-sdk-go/v7/linebot"

type LineTemplateMessage struct {
	AltText  string            `json:"altText"`
	Template *linebot.Template `json:"template"`
}
