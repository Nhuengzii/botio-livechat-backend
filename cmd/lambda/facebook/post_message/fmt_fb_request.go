package main

import (
	"github.com/Nhuengzii/botio-livechat-backend/livechat/api/postmessage"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/external_api/facebook/postfbmessage"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"
)

func fmtFbRequest(req *postmessage.Request, pageID string, psid string) (*postfbmessage.SendingMessage, error) {
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
		switch stdAttachment { // supported post type
		case stdmessage.AttachmentTypeImage:
		case stdmessage.AttachmentTypeVideo:
		case stdmessage.AttachmentTypeAudio:
		case stdmessage.AttachmentTypeFile:
			// add more supported type here
		default:
			return nil, errAttachmentTypeNotSupported
		}
		fbRequest = postfbmessage.SendingMessage{
			Recipient: postfbmessage.Recipient{
				Id: psid,
			},
			MessagingType: "RESPONSE",
			Message: postfbmessage.MessageAttachment{
				Attachment: postfbmessage.AttachmentFacebookRequest{
					AttachmentType: req.Attachment.AttachmentType,
					Payload: postfbmessage.AttachmentFacebookPayload{
						Src:        req.Attachment.Payload.Src,
						IsReusable: true,
					},
				},
			},
		}
	}

	return &fbRequest, nil
}
