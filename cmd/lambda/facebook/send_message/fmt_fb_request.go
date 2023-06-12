package main

import (
	"github.com/Nhuengzii/botio-livechat-backend/livechat/api/sendmessage"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/external_api/facebook/sendfbmessage"
)

func fmtFbRequest(req *sendmessage.Request, pageID string, psid string) *sendfbmessage.SendingMessage {
	var fbRequest sendfbmessage.SendingMessage
	if req.Message != "" {
		fbRequest = sendfbmessage.SendingMessage{
			Recipient: sendfbmessage.Recipient{
				Id: psid,
			},
			MessagingType: "RESPONSE",
			Message: sendfbmessage.MessageText{
				Text: req.Message,
			},
		}
	} else {
		fbRequest = sendfbmessage.SendingMessage{
			Recipient: sendfbmessage.Recipient{
				Id: psid,
			},
			MessagingType: "RESPONSE",
			Message: sendfbmessage.MessageAttachment{
				Attachment: sendfbmessage.AttachmentFacebookRequest{
					AttachmentType: req.Attachment.AttachmentType,
					Payload: sendfbmessage.AttachmentFacebookPayload{
						Src:        req.Attachment.Payload.Src,
						IsReusable: true,
					},
				},
			},
		}
	}

	return &fbRequest
}
