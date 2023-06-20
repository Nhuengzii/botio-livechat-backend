package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/api/postmessage"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var errUnsupportedAttachmentType = errors.New("unsupported attachment type")

func (c *config) handlePostMessageRequest(ctx context.Context, shopID string, pageID string, conversationID string, bot *linebot.Client, requestBody postmessage.Request) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("handlePostMessageRequest: %w", err)
		}
	}()

	attachments := []stdmessage.Attachment{}
	// text message
	if requestBody.Message != "" {
		_, err = bot.PushMessage(conversationID, toLineTextMessage(requestBody)).Do()
		if err != nil {
			return err
		}
	} else {
		// attachment
		switch stdmessage.AttachmentType(requestBody.Attachment.AttachmentType) {
		case stdmessage.AttachmentTypeImage:
			_, err = bot.PushMessage(conversationID, toLineImageMessage(requestBody)).Do()
			if err != nil {
				return err
			}
			attachment := stdmessage.Attachment{
				AttachmentType: stdmessage.AttachmentType(requestBody.Attachment.AttachmentType),
				Payload:        stdmessage.Payload{Src: requestBody.Attachment.Payload.Src},
			}
			attachments = append(attachments, attachment)
		case stdmessage.AttachmentTypeVideo:
			_, err = bot.PushMessage(conversationID, toLineVideoMessage(requestBody)).Do()
			if err != nil {
				return err
			}
			attachment := stdmessage.Attachment{
				AttachmentType: stdmessage.AttachmentType(requestBody.Attachment.AttachmentType),
				Payload:        stdmessage.Payload{Src: requestBody.Attachment.Payload.Src},
			}
			attachments = append(attachments, attachment)
		case stdmessage.AttachmentTypeAudio:
			_, err = bot.PushMessage(conversationID, toLineAudioMessage(requestBody)).Do()
			if err != nil {
				return err
			}
			attachment := stdmessage.Attachment{
				AttachmentType: stdmessage.AttachmentType(requestBody.Attachment.AttachmentType),
				Payload:        stdmessage.Payload{Src: requestBody.Attachment.Payload.Src},
			}
			attachments = append(attachments, attachment)
		case stdmessage.AttachmentTypeLineTemplateButtons:
		case stdmessage.AttachmentTypeLineTemplateConfirm:
		case stdmessage.AttachmentTypeLineTemplateCarousel:
		case stdmessage.AttachmentTypeLineTemplateImageCarousel:
		default:
			return fmt.Errorf("%w: %v", errUnsupportedAttachmentType, requestBody.Attachment.AttachmentType)
		}
	}

	stdMessage := stdmessage.StdMessage{
		ShopID:         shopID,
		Platform:       stdmessage.PlatformLine,
		PageID:         pageID,
		ConversationID: conversationID,
		MessageID:      fmt.Sprintf("%s:%s:%s:%s", shopID, pageID, conversationID, strconv.FormatInt(time.Now().UnixMilli(), 10)),
		Timestamp:      time.Now().UnixMilli(),
		Source: stdmessage.Source{
			UserID:   pageID,
			UserType: stdmessage.UserTypeAdmin,
		},
		Message:     requestBody.Message,
		Attachments: attachments,
		ReplyTo:     nil,
	}
	err = c.dbClient.UpdateConversationOnNewMessage(ctx, &stdMessage)
	if err != nil {
		return err
	}
	err = c.dbClient.InsertMessage(ctx, &stdMessage)
	if err != nil {
		return err
	}
	return nil
}
