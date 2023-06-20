package main

import (
	"github.com/Nhuengzii/botio-livechat-backend/livechat/api/postmessage"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func toLineTextMessage(requestBody postmessage.Request) *linebot.TextMessage {
	return linebot.NewTextMessage(requestBody.Message)
}

func toLineImageMessage(requestBody postmessage.Request) *linebot.ImageMessage {
	return linebot.NewImageMessage(requestBody.Attachment.Payload.Src, requestBody.Attachment.Payload.Src)
}

func toLineVideoMessage(requestBody postmessage.Request) *linebot.VideoMessage {
	return linebot.NewVideoMessage(requestBody.Attachment.Payload.Src, requestBody.Attachment.Payload.Src)
}

func toLineAudioMessage(requestBody postmessage.Request) *linebot.AudioMessage {
	return linebot.NewAudioMessage(requestBody.Attachment.Payload.Src, 30) // how to get the duration?
}
