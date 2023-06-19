package main

import (
	"github.com/Nhuengzii/botio-livechat-backend/livechat/api/postmessage"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/external_api/facebook/postfbmessage"
)

func fmtBasicPayload(payload postmessage.Payload) (*postfbmessage.AttachmentFacebookPayload, error) {
	if payload.Src == "" {
		return nil, errNoSrcFoundForBasicPayload
	}
	return &postfbmessage.AttachmentFacebookPayload{
		Src:        payload.Src,
		IsReusable: true,
	}, nil
}

// func fmtGenericTemplatePayload(payload postmessage.Payload) (*postfbmessage.AttachmentFacebookPayload, error) {
//   if len(payload.FBTemplateGeneric) == 0 {
// return nil,errNoPayloadFoundForTemplatePayload
// }
//
// }
