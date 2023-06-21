package main

import (
	"github.com/Nhuengzii/botio-livechat-backend/livechat/api/postmessage"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/external_api/facebook/reqfbsendmessage"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"
)

func fmtFbRequest(req *postmessage.Request, psid string) (_ *reqfbsendmessage.SendingMessage, err error) {
	var fbRequest reqfbsendmessage.SendingMessage

	if req.Message != "" {
		fbRequest = reqfbsendmessage.SendingMessage{
			Recipient: reqfbsendmessage.Recipient{
				Id: psid,
			},
			MessagingType: "RESPONSE",
			Message: reqfbsendmessage.MessageText{
				Text: req.Message,
			},
		}
	} else {
		stdAttachment := stdmessage.AttachmentType(req.Attachment.AttachmentType)
		var payload *reqfbsendmessage.AttachmentFacebookPayload
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
		case stdmessage.AttachmentTypeFBTemplateGeneric:
			attachmentType = attachmentTypeTemplate
			payload, err = fmtGenericTemplatePayload(req.Attachment.Payload)
			// add more supported type here
		default:
			return nil, errAttachmentTypeNotSupported
		}
		if err != nil {
			return nil, err
		}
		fbRequest = reqfbsendmessage.SendingMessage{
			Recipient: reqfbsendmessage.Recipient{
				Id: psid,
			},
			MessagingType: "RESPONSE",
			Message: reqfbsendmessage.MessageAttachment{
				Attachment: reqfbsendmessage.AttachmentFacebookRequest{
					AttachmentType: attachmentType,
					Payload:        *payload,
				},
			},
		}
	}

	return &fbRequest, nil
}
