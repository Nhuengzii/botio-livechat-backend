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
			} else {
				attachments, err = c.fmtTemplateAttachments(attachment)
			}
			if err != nil {
				return nil, err
			}
		}
	}

	return attachments, nil
}

func (c *config) fmtBasicAttachments(attachment Attachment) ([]stdmessage.Attachment, error) {
	attachments := []stdmessage.Attachment{}
	location, err := getAndUploadMessageContent(c.uploader, attachment.Payload.Src)
	if err != nil {
		return nil, err
	}
	attachments = append(attachments, stdmessage.Attachment{
		AttachmentType: stdmessage.AttachmentType(attachment.AttachmentType),
		Payload:        stdmessage.Payload{Src: location},
	})

	return attachments, nil
}

func (c *config) fmtTemplateAttachments(attachment Attachment) ([]stdmessage.Attachment, error) {
	attachments := []stdmessage.Attachment{}

	var attachmentType stdmessage.AttachmentType
	if attachment.Payload.TemplateType == "button" {
		attachmentType = stdmessage.AttachmentTypeFBTemplateButton
	} else if attachment.Payload.TemplateType == "coupon" {
		attachmentType = stdmessage.AttachmentTypeFBTemplateCoupon
		if attachment.Payload.ImageURL != "" {
			location, err := getAndUploadMessageContent(c.uploader, attachment.Payload.ImageURL)
			if err != nil {
				return nil, err
			}
			attachment.Payload.ImageURL = location
		}
	} else if attachment.Payload.TemplateType == "customer_feedback" {
		attachmentType = stdmessage.AttachmentTypeFBTemplateCustomerFeedback
	} else if attachment.Payload.TemplateType == "generic" {
		attachmentType = stdmessage.AttachmentTypeFBTemplateGeneric
		for index, element := range attachment.Payload.Elements {
			if element.ImageURL != "" {
				location, err := getAndUploadMessageContent(c.uploader, element.ImageURL)
				if err != nil {
					return nil, err
				}
				attachment.Payload.Elements[index].ImageURL = location
			}
		}
	} else if attachment.Payload.TemplateType == "media" {
		attachmentType = stdmessage.AttachmentTypeFBTemplateMedia
	} else if attachment.Payload.TemplateType == "product" {
		attachmentType = stdmessage.AttachmentTypeFBTemplateProduct
	} else if attachment.Payload.TemplateType == "receipt" {
		attachmentType = stdmessage.AttachmentTypeFBTemplateReceipt
		for index, element := range attachment.Payload.Elements {
			if element.ImageURL != "" {
				location, err := getAndUploadMessageContent(c.uploader, element.ImageURL)
				if err != nil {
					return nil, err
				}
				attachment.Payload.Elements[index].ImageURL = location
			}
		}
	} else if attachment.Payload.TemplateType == "customer_information" {
		attachmentType = stdmessage.AttachmentTypeFBTemplateStructuredInformation
	} else {
		return nil, errUnknownTemplateType
	}

	jsonByte, err := json.Marshal(attachment.Payload)
	if err != nil {
		return nil, err
	}
	attachments = append(attachments, stdmessage.Attachment{
		AttachmentType: attachmentType,
		Payload:        stdmessage.Payload{Src: string(jsonByte)},
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
