package main

import "github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"

func fmtBasicAttachments(basicPayload BasicPayload, attachmentType string, jsonBytePayload []byte) ([]stdmessage.Attachment, error) {
	attachments := []stdmessage.Attachment{}
	attachments = append(attachments, stdmessage.Attachment{
		AttachmentType: stdmessage.AttachmentType(attachmentType),
		Payload:        stdmessage.Payload{Src: basicPayload.Src},
	})

	return attachments, nil
}

func fmtTemplateAttachments(templatePayload TemplatePayload, jsonBytePayload []byte) ([]stdmessage.Attachment, error) {
	attachments := []stdmessage.Attachment{}
	var attachmentType stdmessage.AttachmentType

	if templatePayload.TemplateType == "button" {
		attachmentType = stdmessage.AttachmentTypeFBTemplateButton
	} else if templatePayload.TemplateType == "coupon" {
		attachmentType = stdmessage.AttachmentTypeFBTemplateCoupon
	} else if templatePayload.TemplateType == "customer_feedback" {
		attachmentType = stdmessage.AttachmentTypeFBTemplateCustomerFeedback
	} else if templatePayload.TemplateType == "generic" {
		attachmentType = stdmessage.AttachmentTypeFBTemplateGeneric
	} else if templatePayload.TemplateType == "media" {
		attachmentType = stdmessage.AttachmentTypeFBTemplateMedia
	} else if templatePayload.TemplateType == "product" {
		attachmentType = stdmessage.AttachmentTypeFBTemplateProduct
	} else if templatePayload.TemplateType == "receipt" {
		attachmentType = stdmessage.AttachmentTypeFBTemplateReceipt
	} else if templatePayload.TemplateType == "customer_information" {
		attachmentType = stdmessage.AttachmentTypeFBTemplateStructuredInformation
	} else {
		return nil, errUnknownTemplateType
	}
	attachments = append(attachments, stdmessage.Attachment{
		AttachmentType: attachmentType,
		Payload:        stdmessage.Payload{Src: string(jsonBytePayload)},
	})
	return attachments, nil
}
