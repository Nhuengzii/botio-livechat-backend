package main

import (
	"github.com/Nhuengzii/botio-livechat-backend/livechat/api/postmessage"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func toLineTextMessage(req postmessage.Request) *linebot.TextMessage {
	return linebot.NewTextMessage(req.Message)
}

func toLineImageMessage(req postmessage.Request) *linebot.ImageMessage {
	return linebot.NewImageMessage(req.Attachment.Payload.Src, req.Attachment.Payload.Src)
}

func toLineVideoMessage(req postmessage.Request) *linebot.VideoMessage {
	return linebot.NewVideoMessage(req.Attachment.Payload.Src, req.Attachment.Payload.Src)
}

func toLineAudioMessage(req postmessage.Request) *linebot.AudioMessage {
	return linebot.NewAudioMessage(req.Attachment.Payload.Src, 30) // how to get the duration?
}

func toLineButtonsTemplateMessage(req postmessage.Request) *linebot.TemplateMessage {
	altText := req.Attachment.Payload.LineTemplateButtons.AltText

	defaultActionLabel := req.Attachment.Payload.LineTemplateButtons.DefaultAction.Label
	defaultActionURI := req.Attachment.Payload.LineTemplateButtons.DefaultAction.URI
	defaultAction := linebot.NewURIAction(defaultActionLabel, defaultActionURI)

	actions := []linebot.TemplateAction{}
	for _, action := range req.Attachment.Payload.LineTemplateButtons.Actions {
		actions = append(actions, linebot.NewURIAction(action.Label, action.URI))
	}

	buttonsTemplate := &linebot.ButtonsTemplate{
		ThumbnailImageURL: req.Attachment.Payload.LineTemplateButtons.ThumbnailImageURL,
		Title:             req.Attachment.Payload.LineTemplateButtons.Title,
		Text:              req.Attachment.Payload.LineTemplateButtons.Text,
		DefaultAction:     defaultAction,
		Actions:           actions,
	}

	return linebot.NewTemplateMessage(altText, buttonsTemplate)
}

func toLineConfirmTemplateMessage(req postmessage.Request) *linebot.TemplateMessage {
	altText := req.Attachment.Payload.LineTemplateConfirm.AltText

	actions := []linebot.TemplateAction{}
	for _, action := range req.Attachment.Payload.LineTemplateConfirm.Actions {
		actions = append(actions, linebot.NewURIAction(action.Label, action.URI))
	}

	confirmTemplate := &linebot.ConfirmTemplate{
		Text:    req.Attachment.Payload.LineTemplateConfirm.Text,
		Actions: actions,
	}

	return linebot.NewTemplateMessage(altText, confirmTemplate)
}

func toLineCarouselTemplateMessage(req postmessage.Request) *linebot.TemplateMessage {
	altText := req.Attachment.Payload.LineTemplateCarousel.AltText

	columns := []*linebot.CarouselColumn{}
	for _, column := range req.Attachment.Payload.LineTemplateCarousel.Columns {
		defaultActionLabel := column.DefaultAction.Label
		defaultActionURI := column.DefaultAction.URI
		defaultAction := linebot.NewURIAction(defaultActionLabel, defaultActionURI)

		actions := []linebot.TemplateAction{}
		for _, action := range column.Actions {
			actions = append(actions, linebot.NewURIAction(action.Label, action.URI))
		}

		columns = append(columns, &linebot.CarouselColumn{
			ThumbnailImageURL: column.ThumbnailImageURL,
			Title:             column.Title,
			Text:              column.Text,
			DefaultAction:     defaultAction,
			Actions:           actions,
		})
	}

	carouselTemplate := &linebot.CarouselTemplate{
		Columns: columns,
	}

	return linebot.NewTemplateMessage(altText, carouselTemplate)
}

func toLineImageCarouselTemplateMessage(req postmessage.Request) *linebot.TemplateMessage {
	altText := req.Attachment.Payload.LineTemplateImageCarousel.AltText

	columns := []*linebot.ImageCarouselColumn{}
	for _, column := range req.Attachment.Payload.LineTemplateImageCarousel.Columns {
		columns = append(columns, &linebot.ImageCarouselColumn{
			ImageURL: column.ImageURL,
			Action:   linebot.NewURIAction(column.Action.Label, column.Action.URI),
		})
	}

	imageCarouselTemplate := &linebot.ImageCarouselTemplate{
		Columns: columns,
	}

	return linebot.NewTemplateMessage(altText, imageCarouselTemplate)
}
