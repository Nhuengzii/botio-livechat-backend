package main

import (
	"fmt"
	"net/http"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/storage/amazons3"
)

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

func getAndUploadMessageContent(uploader amazons3.Uploader, src string) (_ string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("getAndUploadMessageContent: %w", err)
		}
	}()
	resp, err := http.Get(src)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errBadImageURL
	}

	location, err := uploader.UploadFile(resp.Body)
	if err != nil {
		return "", err
	}
	return location, nil
}
