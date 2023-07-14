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

	actions := []linebot.TemplateAction{}
	for _, action := range req.Attachment.Payload.LineTemplateButtons.Actions {
		actions = append(actions, linebot.NewURIAction(action.Label, action.URI))
	}

	buttonsTemplate := &linebot.ButtonsTemplate{
		ThumbnailImageURL: req.Attachment.Payload.LineTemplateButtons.ThumbnailImageURL,
		Title:             req.Attachment.Payload.LineTemplateButtons.Title,
		Text:              req.Attachment.Payload.LineTemplateButtons.Text,
		Actions:           actions,
	}

	if req.Attachment.Payload.LineTemplateButtons.DefaultAction != nil {
		defaultAction := linebot.NewURIAction(
			req.Attachment.Payload.LineTemplateButtons.DefaultAction.Label,
			req.Attachment.Payload.LineTemplateButtons.DefaultAction.URI,
		)
		buttonsTemplate.WithDefaultAction(defaultAction)
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
	for _, reqCarouselColumn := range req.Attachment.Payload.LineTemplateCarousel.Columns {
		actions := []linebot.TemplateAction{}
		for _, action := range reqCarouselColumn.Actions {
			actions = append(actions, linebot.NewURIAction(action.Label, action.URI))
		}

		column := &linebot.CarouselColumn{
			ThumbnailImageURL: reqCarouselColumn.ThumbnailImageURL,
			Title:             reqCarouselColumn.Title,
			Text:              reqCarouselColumn.Text,
			Actions:           actions,
		}

		if reqCarouselColumn.DefaultAction != nil {
			defaultAction := linebot.NewURIAction(
				reqCarouselColumn.DefaultAction.Label,
				reqCarouselColumn.DefaultAction.URI,
			)
			column.WithDefaultAction(defaultAction)
		}

		columns = append(columns, column)
	}

	carouselTemplate := &linebot.CarouselTemplate{
		Columns: columns,
	}

	return linebot.NewTemplateMessage(altText, carouselTemplate)
}

func toLineImageCarouselTemplateMessage(req postmessage.Request) *linebot.TemplateMessage {
	altText := req.Attachment.Payload.LineTemplateImageCarousel.AltText

	columns := []*linebot.ImageCarouselColumn{}
	for _, reqImageCarouselColumn := range req.Attachment.Payload.LineTemplateImageCarousel.Columns {
		columns = append(columns, &linebot.ImageCarouselColumn{
			ImageURL: reqImageCarouselColumn.ImageURL,
			Action:   linebot.NewURIAction(reqImageCarouselColumn.Action.Label, reqImageCarouselColumn.Action.URI),
		})
	}

	imageCarouselTemplate := &linebot.ImageCarouselTemplate{
		Columns: columns,
	}

	return linebot.NewTemplateMessage(altText, imageCarouselTemplate)
}
