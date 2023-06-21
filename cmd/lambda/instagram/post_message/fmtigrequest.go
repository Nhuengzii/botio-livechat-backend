package main

import (
	"github.com/Nhuengzii/botio-livechat-backend/livechat/api/postmessage"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/external_api/instagram/reqigsendmessage"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"
)

func fmtIgRequest(req *postmessage.Request, psid string) (_ *reqigsendmessage.SendingMessage, err error) {
	var igRequest reqigsendmessage.SendingMessage

	if req.Message != "" {
		igRequest = reqigsendmessage.SendingMessage{
			Recipient: reqigsendmessage.Recipient{
				Id: psid,
			},
			Message: reqigsendmessage.MessageText{
				Text: req.Message,
			},
		}
	} else {
		stdAttachment := stdmessage.AttachmentType(req.Attachment.AttachmentType)
		var payload *reqigsendmessage.AttachmentInstagramPayload
		var attachmentType string
		switch stdAttachment { // supported post type
		case stdmessage.AttachmentTypeImage:
			attachmentType = req.Attachment.AttachmentType
			payload, err = fmtBasicPayload(req.Attachment.Payload)
		case stdmessage.AttachmentTypeVideo:
			attachmentType = req.Attachment.AttachmentType
			payload, err = fmtBasicPayload(req.Attachment.Payload)
		case stdmessage.AttachmentTypeAudio:
			attachmentType = req.Attachment.AttachmentType
			payload, err = fmtBasicPayload(req.Attachment.Payload)
		case stdmessage.AttachmentTypeFile:
			attachmentType = req.Attachment.AttachmentType
			payload, err = fmtBasicPayload(req.Attachment.Payload)
		case stdmessage.AttachmentTypeIGTemplateGeneric:
			attachmentType = attachmentTypeTemplate
			payload, err = fmtGenericTemplatePayload(req.Attachment.Payload)
			// add more supported type here
		default:
			return nil, errAttachmentTypeNotSupported
		}
		if err != nil {
			return nil, err
		}
		igRequest = reqigsendmessage.SendingMessage{
			Recipient: reqigsendmessage.Recipient{
				Id: psid,
			},
			Message: reqigsendmessage.MessageAttachment{
				Attachment: reqigsendmessage.AttachmentInstagramRequest{
					AttachmentType: attachmentType,
					Payload:        *payload,
				},
			},
		}
	}

	return &igRequest, nil
}
