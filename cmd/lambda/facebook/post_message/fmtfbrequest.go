package main

import (
	"github.com/Nhuengzii/botio-livechat-backend/livechat/api/postmessage"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/external_api/facebook/postfbmessage"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"
)

func fmtFbRequest(req *postmessage.Request, pageID string, psid string) (_ *postfbmessage.SendingMessage, err error) {
	var fbRequest postfbmessage.SendingMessage

	if req.Message != "" {
		fbRequest = postfbmessage.SendingMessage{
			Recipient: postfbmessage.Recipient{
				Id: psid,
			},
			MessagingType: "RESPONSE",
			Message: postfbmessage.MessageText{
				Text: req.Message,
			},
		}
	} else {
		stdAttachment := stdmessage.AttachmentType(req.Attachment.AttachmentType)
		var payload *postfbmessage.AttachmentFacebookPayload
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
		fbRequest = postfbmessage.SendingMessage{
			Recipient: postfbmessage.Recipient{
				Id: psid,
			},
			MessagingType: "RESPONSE",
			Message: postfbmessage.MessageAttachment{
				Attachment: postfbmessage.AttachmentFacebookRequest{
					AttachmentType: attachmentType,
					Payload:        *payload,
				},
			},
		}
	}

	return &fbRequest, nil
}
