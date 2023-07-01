package main

import (
	"github.com/Nhuengzii/botio-livechat-backend/livechat/api/postmessage"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/external_api/instagram/reqigsendmessage"
)

func fmtBasicPayload(payload postmessage.Payload) (*reqigsendmessage.AttachmentInstagramPayload, error) {
	if payload.Src == "" {
		return nil, errNoSrcFoundForBasicPayload
	}
	return &reqigsendmessage.AttachmentInstagramPayload{
		Src:        payload.Src,
		IsReusable: true,
	}, nil
}

func fmtGenericTemplatePayload(payload postmessage.Payload) (*reqigsendmessage.AttachmentInstagramPayload, error) {
	if len(payload.IGTemplateGeneric) == 0 {
		return nil, errNoPayloadFoundForTemplatePayload
	}
	var genericTemplate []reqigsendmessage.Template // reqigsendmessage.GenericTemplate
	for _, element := range payload.IGTemplateGeneric {
		var buttons []reqigsendmessage.Button
		for _, button := range element.Button {
			buttons = append(buttons, reqigsendmessage.Button{
				Type:  templateButtonURLType,
				URL:   button.URL,
				Title: button.Title,
			})
		}

		var defaultAction *reqigsendmessage.DefaultAction
		if element.DefaultAction != nil {
			defaultAction = &reqigsendmessage.DefaultAction{
				Type: templateButtonURLType,
				URL:  element.DefaultAction.URL,
			}
		}
		genericTemplate = append(genericTemplate, reqigsendmessage.GenericTemplate{
			Title:         element.Title,
			Subtitle:      element.Message,
			DefaultAction: defaultAction,
			Buttons:       buttons,
			ImageURL:      element.Picture,
		})
	}
	return &reqigsendmessage.AttachmentInstagramPayload{
		TemplateType: templateTypeGeneric,
		Elements:     genericTemplate,
	}, nil
}
