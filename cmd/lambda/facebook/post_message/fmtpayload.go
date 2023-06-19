package main

import (
	"github.com/Nhuengzii/botio-livechat-backend/livechat/api/postmessage"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/external_api/facebook/postfbmessage"
)

func fmtBasicPayload(payload postmessage.Payload) (*postfbmessage.AttachmentFacebookPayload, error) {
	if payload.Src == "" {
		return nil, errNoSrcFoundForBasicPayload
	}
	return &postfbmessage.AttachmentFacebookPayload{
		Src:        payload.Src,
		IsReusable: true,
	}, nil
}

func fmtGenericTemplatePayload(payload postmessage.Payload) (*postfbmessage.AttachmentFacebookPayload, error) {
	if len(payload.FBTemplateGeneric) == 0 {
		return nil, errNoPayloadFoundForTemplatePayload
	}
	var genericTemplate []any // postfbmessage.GenericTemplate
	for _, element := range payload.FBTemplateGeneric {
		var buttons []postfbmessage.Button
		for _, button := range element.Button {
			buttons = append(buttons, postfbmessage.Button{
				Type:  templateButtonURLType,
				URL:   button.URL,
				Title: button.Title,
			})
		}
		genericTemplate = append(genericTemplate, postfbmessage.GenericTemplate{
			Title:    element.Title,
			Subtitle: element.Message,
			DefaultAction: postfbmessage.DefaultAction{
				Type: templateButtonURLType,
				URL:  element.DefaultAction.URL,
			},
			Buttons:  buttons,
			ImageURL: element.Picture,
		})
	}

	return &postfbmessage.AttachmentFacebookPayload{
		TemplateType: templateTypeGeneric,
		Elements:     genericTemplate,
	}, nil
}
