package main

import (
	"github.com/Nhuengzii/botio-livechat-backend/livechat/api/request/postmessagreq"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/external/fbrequest"
)

func fmtFbRequest(req *postmessagreq.Request, pageID string, psid string) *fbrequest.FBSendMsgRequest {
	var fbRequest fbrequest.FBSendMsgRequest
	if req.Message != "" {
		fbRequest = fbrequest.FBSendMsgRequest{
			Recipient: fbrequest.Recipient{
				Id: psid,
			},
			MessagingType: "RESPONSE",
			Message: fbrequest.MessageText{
				Text: req.Message,
			},
		}
	} else {
		fbRequest = fbrequest.FBSendMsgRequest{
			Recipient: fbrequest.Recipient{
				Id: psid,
			},
			MessagingType: "RESPONSE",
			Message: fbrequest.MessageAttachment{
				Attachment: fbrequest.AttachmentFacebookRequest{
					AttachmentType: req.Attachment.AttachmentType,
					Payload: fbrequest.AttachmentFacebookPayload{
						Src:        req.Attachment.Payload.Src,
						IsReusable: true,
					},
				},
			},
		}
	}

	return &fbRequest
}
