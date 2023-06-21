package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/api/postmessage"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var errUnsupportedAttachmentType = errors.New("unsupported attachment type")

func (c *config) handlePostMessageRequest(ctx context.Context, shopID string, pageID string, conversationID string, bot *linebot.Client, req postmessage.Request) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("handlePostMessageRequest: %w", err)
		}
	}()

	// check conversation exists
	err = c.dbClient.CheckConversationExists(ctx, conversationID)
	if err != nil {
		return err
	}

	attachments := []stdmessage.Attachment{}
	// text message
	if req.Message != "" {
		_, err = bot.PushMessage(conversationID, toLineTextMessage(req)).Do()
		if err != nil {
			return err
		}
	} else {
		// attachment
		switch stdmessage.AttachmentType(req.Attachment.AttachmentType) {
		case stdmessage.AttachmentTypeImage:
			_, err = bot.PushMessage(conversationID, toLineImageMessage(req)).Do()
			if err != nil {
				return err
			}
			attachment := stdmessage.Attachment{
				AttachmentType: stdmessage.AttachmentType(req.Attachment.AttachmentType),
				Payload:        stdmessage.Payload{Src: req.Attachment.Payload.Src},
			}
			attachments = append(attachments, attachment)
		case stdmessage.AttachmentTypeVideo:
			_, err = bot.PushMessage(conversationID, toLineVideoMessage(req)).Do()
			if err != nil {
				return err
			}
			attachment := stdmessage.Attachment{
				AttachmentType: stdmessage.AttachmentType(req.Attachment.AttachmentType),
				Payload:        stdmessage.Payload{Src: req.Attachment.Payload.Src},
			}
			attachments = append(attachments, attachment)
		case stdmessage.AttachmentTypeAudio:
			_, err = bot.PushMessage(conversationID, toLineAudioMessage(req)).Do()
			if err != nil {
				return err
			}
			attachment := stdmessage.Attachment{
				AttachmentType: stdmessage.AttachmentType(req.Attachment.AttachmentType),
				Payload:        stdmessage.Payload{Src: req.Attachment.Payload.Src},
			}
			attachments = append(attachments, attachment)
		case stdmessage.AttachmentTypeLineTemplateButtons:
			_, err = bot.PushMessage(conversationID, toLineButtonsTemplateMessage(req)).Do()
			if err != nil {
				return err
			}
			stdMessagePayloadSrcJSON, err := json.Marshal(
				struct {
					AttachmentType stdmessage.AttachmentType       `json:"attachmentType"`
					Payload        postmessage.LineTemplateButtons `json:"payload"`
				}{
					AttachmentType: stdmessage.AttachmentType(req.Attachment.AttachmentType),
					Payload:        req.Attachment.Payload.LineTemplateButtons,
				},
			)
			if err != nil {
				return err
			}
			attachment := stdmessage.Attachment{
				AttachmentType: stdmessage.AttachmentType(req.Attachment.AttachmentType),
				Payload:        stdmessage.Payload{Src: string(stdMessagePayloadSrcJSON)},
			}
			attachments = append(attachments, attachment)
		case stdmessage.AttachmentTypeLineTemplateConfirm:
			_, err = bot.PushMessage(conversationID, toLineConfirmTemplateMessage(req)).Do()
			if err != nil {
				return err
			}
			stdMessagePayloadSrcJSON, err := json.Marshal(
				struct {
					AttachmentType stdmessage.AttachmentType       `json:"attachmentType"`
					Payload        postmessage.LineTemplateConfirm `json:"payload"`
				}{
					AttachmentType: stdmessage.AttachmentType(req.Attachment.AttachmentType),
					Payload:        req.Attachment.Payload.LineTemplateConfirm,
				},
			)
			if err != nil {
				return err
			}
			attachment := stdmessage.Attachment{
				AttachmentType: stdmessage.AttachmentType(req.Attachment.AttachmentType),
				Payload:        stdmessage.Payload{Src: string(stdMessagePayloadSrcJSON)},
			}
			attachments = append(attachments, attachment)
		case stdmessage.AttachmentTypeLineTemplateCarousel:
			_, err = bot.PushMessage(conversationID, toLineCarouselTemplateMessage(req)).Do()
			if err != nil {
				return err
			}
			stdMessagePayloadSrcJSON, err := json.Marshal(
				struct {
					AttachmentType stdmessage.AttachmentType        `json:"attachmentType"`
					Payload        postmessage.LineTemplateCarousel `json:"payload"`
				}{
					AttachmentType: stdmessage.AttachmentType(req.Attachment.AttachmentType),
					Payload:        req.Attachment.Payload.LineTemplateCarousel,
				},
			)
			if err != nil {
				return err
			}
			attachment := stdmessage.Attachment{
				AttachmentType: stdmessage.AttachmentType(req.Attachment.AttachmentType),
				Payload:        stdmessage.Payload{Src: string(stdMessagePayloadSrcJSON)},
			}
			attachments = append(attachments, attachment)
		case stdmessage.AttachmentTypeLineTemplateImageCarousel:
			_, err = bot.PushMessage(conversationID, toLineImageCarouselTemplateMessage(req)).Do()
			if err != nil {
				return err
			}
			stdMessagePayloadSrcJSON, err := json.Marshal(
				struct {
					AttachmentType stdmessage.AttachmentType             `json:"attachmentType"`
					Payload        postmessage.LineTemplateImageCarousel `json:"payload"`
				}{
					AttachmentType: stdmessage.AttachmentType(req.Attachment.AttachmentType),
					Payload:        req.Attachment.Payload.LineTemplateImageCarousel,
				},
			)
			if err != nil {
				return err
			}
			attachment := stdmessage.Attachment{
				AttachmentType: stdmessage.AttachmentType(req.Attachment.AttachmentType),
				Payload:        stdmessage.Payload{Src: string(stdMessagePayloadSrcJSON)},
			}
			attachments = append(attachments, attachment)
		default:
			return fmt.Errorf("%w: %v", errUnsupportedAttachmentType, req.Attachment.AttachmentType)
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
		Message:     req.Message,
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
