package main

import (
	"github.com/Nhuengzii/botio-livechat-backend/livechat/api/postmessage"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/discord"
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

func toLineTemplateMessage(requestBody postmessage.Request) *linebot.TemplateMessage {
	discord.Log("https://discord.com/api/webhooks/1109019632339267584/C26EwyFL2Njn7iLX9VDIto4uF_5C7Qqm3aKuUthHKbJYGLoNM_394GddBbW5gqYPP6Ei", "toLineTemplateMessage")
	templateMap := requestBody.Attachment.Payload.LineTemplateMessage.Template.(map[string]interface{})
	switch templateMap["type"].(linebot.TemplateType) {
	case linebot.TemplateTypeButtons:
		return linebot.NewTemplateMessage(requestBody.Attachment.Payload.LineTemplateMessage.AltText, toLineButtonsTemplate(requestBody))
	case linebot.TemplateTypeConfirm:
		discord.Log("https://discord.com/api/webhooks/1109019632339267584/C26EwyFL2Njn7iLX9VDIto4uF_5C7Qqm3aKuUthHKbJYGLoNM_394GddBbW5gqYPP6Ei", "case line confirm template")
		return linebot.NewTemplateMessage(requestBody.Attachment.Payload.LineTemplateMessage.AltText, toLineConfirmTemplate(requestBody))
	case linebot.TemplateTypeCarousel:
		return linebot.NewTemplateMessage(requestBody.Attachment.Payload.LineTemplateMessage.AltText, toLineCarouselTemplate(requestBody))
	case linebot.TemplateTypeImageCarousel:
		return linebot.NewTemplateMessage(requestBody.Attachment.Payload.LineTemplateMessage.AltText, toLineImageCarouselTemplate(requestBody))
	default:
		discord.Log("https://discord.com/api/webhooks/1109019632339267584/C26EwyFL2Njn7iLX9VDIto4uF_5C7Qqm3aKuUthHKbJYGLoNM_394GddBbW5gqYPP6Ei", "case default")
		return nil
	}
}

func toLineButtonsTemplate(requestBody postmessage.Request) *linebot.ButtonsTemplate {
	templateMap := requestBody.Attachment.Payload.LineTemplateMessage.Template.(map[string]interface{})
	thumbnailImageURL := templateMap["thumbnailImageUrl"].(string)
	imageAspectRatio := templateMap["imageAspectRatio"].(linebot.ImageAspectRatioType)
	imageSize := templateMap["imageSize"].(linebot.ImageSizeType)
	imageBackgroundColor := templateMap["imageBackgroundColor"].(string)
	title := templateMap["title"].(string)
	text := templateMap["text"].(string)
	actions := templateMap["actions"].([]linebot.TemplateAction)
	defaultAction := templateMap["defaultAction"].(linebot.TemplateAction)
	return &linebot.ButtonsTemplate{
		ThumbnailImageURL:    thumbnailImageURL,
		ImageAspectRatio:     imageAspectRatio,
		ImageSize:            imageSize,
		ImageBackgroundColor: imageBackgroundColor,
		Title:                title,
		Text:                 text,
		Actions:              actions,
		DefaultAction:        defaultAction,
	}
}

func toLineConfirmTemplate(requestBody postmessage.Request) *linebot.ConfirmTemplate {
	discord.Log("https://discord.com/api/webhooks/1109019632339267584/C26EwyFL2Njn7iLX9VDIto4uF_5C7Qqm3aKuUthHKbJYGLoNM_394GddBbW5gqYPP6Ei", "toLineConfirmTemplate")
	templateMap := requestBody.Attachment.Payload.LineTemplateMessage.Template.(map[string]interface{})
	text := templateMap["text"].(string)
	actions := templateMap["actions"].([]linebot.TemplateAction)
	return &linebot.ConfirmTemplate{
		Text:    text,
		Actions: actions,
	}
}

func toLineCarouselTemplate(requestBody postmessage.Request) *linebot.CarouselTemplate {
	templateMap := requestBody.Attachment.Payload.LineTemplateMessage.Template.(map[string]interface{})
	columns := templateMap["columns"].([]*linebot.CarouselColumn)
	imageAspectRatio := templateMap["imageAspectRatio"].(linebot.ImageAspectRatioType)
	imageSize := templateMap["imageSize"].(linebot.ImageSizeType)
	return &linebot.CarouselTemplate{
		Columns:          columns,
		ImageAspectRatio: imageAspectRatio,
		ImageSize:        imageSize,
	}
}

func toLineImageCarouselTemplate(requestBody postmessage.Request) *linebot.ImageCarouselTemplate {
	templateMap := requestBody.Attachment.Payload.LineTemplateMessage.Template.(map[string]interface{})
	columns := templateMap["columns"].([]*linebot.ImageCarouselColumn)
	return &linebot.ImageCarouselTemplate{
		Columns: columns,
	}
}
