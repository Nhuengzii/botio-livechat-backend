package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/storage/amazons3"
)

func (c *config) fmtAttachment(messaging Messaging) ([]stdmessage.Attachment, error) {
	attachments := []stdmessage.Attachment{}
	if len(messaging.Message.Attachments) > 0 {
		for _, attachment := range messaging.Message.Attachments {
			if attachment.AttachmentType != "template" {
				tempAttachment, err := c.fmtBasicAttachments(attachment)
				if err != nil {
					return nil, err
				}
				attachments = append(attachments, tempAttachment)
			} else {
				tempAttachment, err := c.fmtTemplateAttachments(attachment)
				if err != nil {
					return nil, err
				}
				attachments = append(attachments, tempAttachment)
			}
		}
	}

	return attachments, nil
}

func (c *config) fmtBasicAttachments(attachment Attachment) (stdmessage.Attachment, error) {
	location, err := getAndUploadMessageContent(c.uploader, attachment.Payload.Src)
	if err != nil {
		return stdmessage.Attachment{}, err
	}
	return stdmessage.Attachment{
		AttachmentType: stdmessage.AttachmentType(attachment.AttachmentType),
		Payload:        stdmessage.Payload{Src: location},
	}, nil
}

func (c *config) fmtTemplateAttachments(attachment Attachment) (stdmessage.Attachment, error) {
	var attachmentType stdmessage.AttachmentType
	if attachment.Payload.TemplateType == "button" {
		attachmentType = stdmessage.AttachmentTypeFBTemplateButton
	} else if attachment.Payload.TemplateType == "coupon" {
		attachmentType = stdmessage.AttachmentTypeFBTemplateCoupon
		if attachment.Payload.ImageURL != "" {
			location, err := getAndUploadMessageContent(c.uploader, attachment.Payload.ImageURL)
			if err != nil {
				return stdmessage.Attachment{}, err
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
					return stdmessage.Attachment{}, err
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
					return stdmessage.Attachment{}, err
				}
				attachment.Payload.Elements[index].ImageURL = location
			}
		}
	} else if attachment.Payload.TemplateType == "customer_information" {
		attachmentType = stdmessage.AttachmentTypeFBTemplateStructuredInformation
	} else {
		return stdmessage.Attachment{}, errUnknownTemplateType
	}

	jsonByte, err := json.Marshal(attachment.Payload)
	if err != nil {
		return stdmessage.Attachment{}, err
	}
	return stdmessage.Attachment{
		AttachmentType: attachmentType,
		Payload:        stdmessage.Payload{Src: string(jsonByte)},
	}, nil
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

	body, err := io.ReadAll(resp.Body)
	location, err := uploader.UploadFile(body)
	if err != nil {
		return "", err
	}
	return location, nil
}
