package main

import (
	"github.com/Nhuengzii/botio-livechat-backend/livechat/api/postmessage"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/external_api/facebook/reqfbsendmessage"
)

func fmtBasicPayload(payload postmessage.Payload) (*reqfbsendmessage.AttachmentFacebookPayload, error) {
	if payload.Src == "" {
		return nil, errNoSrcFoundForBasicPayload
	}
	return &reqfbsendmessage.AttachmentFacebookPayload{
		Src:        payload.Src,
		IsReusable: true,
	}, nil
}

func fmtGenericTemplatePayload(payload postmessage.Payload) (*reqfbsendmessage.AttachmentFacebookPayload, error) {
	if len(payload.FBTemplateGeneric) == 0 {
		return nil, errNoPayloadFoundForTemplatePayload
	}
	var genericTemplate []any // reqfbsendmessage.GenericTemplate
	for _, element := range payload.FBTemplateGeneric {
		var buttons []reqfbsendmessage.Button
		for _, button := range element.Button {
			buttons = append(buttons, reqfbsendmessage.Button{
				Type:  templateButtonURLType,
				URL:   button.URL,
				Title: button.Title,
			})
		}
		genericTemplate = append(genericTemplate, reqfbsendmessage.GenericTemplate{
			Title:    element.Title,
			Subtitle: element.Message,
			DefaultAction: reqfbsendmessage.DefaultAction{
				Type: templateButtonURLType,
				URL:  element.DefaultAction.URL,
			},
			Buttons:  buttons,
			ImageURL: element.Picture,
		})
	}

	return &reqfbsendmessage.AttachmentFacebookPayload{
		TemplateType: templateTypeGeneric,
		Elements:     genericTemplate,
	}, nil
}
