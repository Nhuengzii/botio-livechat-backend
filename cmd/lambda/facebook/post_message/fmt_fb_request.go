package main

import (
	"github.com/Nhuengzii/botio-livechat-backend/livechat/api/postmessage"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/external_api/facebook/postfbmessage"
)

func fmtFbRequest(req *postmessage.Request, pageID string, psid string) *postfbmessage.SendingMessage {
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

	return &fbRequest
}
