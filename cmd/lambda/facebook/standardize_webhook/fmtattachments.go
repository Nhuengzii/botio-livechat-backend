package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/storage/amazons3"
)

func (c *config) fmtAttachment(messaging Messaging) ([]stdmessage.Attachment, error) {
	var err error
	attachments := []stdmessage.Attachment{}
	if len(messaging.Message.Attachments) > 0 {
		for _, attachment := range messaging.Message.Attachments {
			if attachment.AttachmentType != "template" {
				attachments, err = c.fmtBasicAttachments(attachment)
				if err != nil {
					return nil, err
				}
			} else {
				jsonByte, err := json.Marshal(attachment.Payload) // actual payload
				if err != nil {
					return nil, err
				}
				var templatePayload TemplatePayload
				err = json.Unmarshal(jsonByte, &templatePayload)
				if err != nil {
					return nil, err
				}
				attachments, err = c.fmtTemplateAttachments(templatePayload, jsonByte)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return attachments, nil
}

func (c *config) fmtBasicAttachments(attachment Attachment) ([]stdmessage.Attachment, error) {
	jsonByte, err := json.Marshal(attachment.Payload) // actual payload
	if err != nil {
		return nil, err
	}
	var basicPayload BasicPayload
	err = json.Unmarshal([]byte(jsonByte), &basicPayload)
	if err != nil {
		return nil, err
	}

	attachments := []stdmessage.Attachment{}
	location, err := getAndUploadMessageContent(c.uploader, basicPayload.Src)
	if err != nil {
		return nil, err
	}
	attachments = append(attachments, stdmessage.Attachment{
		AttachmentType: stdmessage.AttachmentType(attachment.AttachmentType),
		Payload:        stdmessage.Payload{Src: location},
	})

	return attachments, nil
}

func (c *config) fmtTemplateAttachments(templatePayload TemplatePayload, jsonBytePayload []byte) ([]stdmessage.Attachment, error) {
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
