package main

import (
	"github.com/Nhuengzii/botio-livechat-backend/livechat/api/postmessage"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/external_api/facebook"
)

func fmtFbRequest(req *postmessage.Request, pageID string, psid string) *facebook.FBSendMsgRequest {
	var fbRequest facebook.FBSendMsgRequest
	if req.Message != "" {
		fbRequest = facebook.FBSendMsgRequest{
			Recipient: facebook.Recipient{
				Id: psid,
			},
			MessagingType: "RESPONSE",
			Message: facebook.MessageText{
				Text: req.Message,
			},
		}
	} else {
		fbRequest = facebook.FBSendMsgRequest{
			Recipient: facebook.Recipient{
				Id: psid,
			},
			MessagingType: "RESPONSE",
			Message: facebook.MessageAttachment{
				Attachment: facebook.AttachmentFacebookRequest{
					AttachmentType: req.Attachment.AttachmentType,
					Payload: facebook.AttachmentFacebookPayload{
						Src:        req.Attachment.Payload.Src,
						IsReusable: true,
					},
				},
			},
		}
	}

	return &fbRequest
}
