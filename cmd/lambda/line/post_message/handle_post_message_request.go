package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/api/postmessage"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func (c *config) handlePostMessageRequest(ctx context.Context, shopID string, pageID string, conversationID string, bot *linebot.Client, requestBody postmessage.Request) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("handlePostMessageRequest: %w", err)
		}
	}()
	if requestBody.Message != "" {
		_, err = bot.PushMessage(conversationID, linebot.NewTextMessage(requestBody.Message)).Do()
		if err != nil {
			return err
		}
	} else {
		switch requestBody.Attachment.AttachmentType {
		case stdmessage.AttachmentTypeImage:
			_, err = bot.PushMessage(conversationID, linebot.NewImageMessage(requestBody.Attachment.Payload.Src, requestBody.Attachment.Payload.Src)).Do()
			if err != nil {
				return err
			}
		case stdmessage.AttachmentTypeVideo:
			_, err = bot.PushMessage(conversationID, linebot.NewVideoMessage(requestBody.Attachment.Payload.Src, requestBody.Attachment.Payload.Src)).Do()
			if err != nil {
				return err
			}
		case stdmessage.AttachmentTypeAudio:
			_, err = bot.PushMessage(conversationID, linebot.NewAudioMessage(requestBody.Attachment.Payload.Src, 30)).Do() // where to get duration?
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("handlePostMessageRequest: unsuported attachment type: %v", requestBody.Attachment.AttachmentType)
		}
	}
	attachment := stdmessage.Attachment{}
	if requestBody.Message == "" {
		attachment = stdmessage.Attachment{
			AttachmentType: requestBody.Attachment.AttachmentType,
			Payload:        stdmessage.Payload{Src: requestBody.Attachment.Payload.Src},
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
		Attachments: []stdmessage.Attachment{attachment},
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
