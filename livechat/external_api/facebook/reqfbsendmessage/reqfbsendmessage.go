package reqfbsendmessage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func SendMessage(accessToken string, message SendingMessage, pageID string) (_ *SendingMessageResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("reqfbsendmessage.SendMessage: %w", err)
		}
	}()
	url := fmt.Sprintf("https://graph.facebook.com/v16.0/%v/messages?access_token=%v", pageID, accessToken)
	facebookReqBody, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(facebookReqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	now := time.Now().UnixMilli()
	var response SendingMessageResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	response.Timestamp = now
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// facebook api response
type SendingMessageResponse struct {
	RecipientID string `json:"recipient_id"`
	MessageID   string `json:"message_id"`
	Timestamp   int64  `json:"timestamp"`
}

// facebook api send message top level struct
type SendingMessage struct {
	Recipient     Recipient `json:"recipient"`
	Message       Message   `json:"message"`
	MessagingType string    `json:"messaging_type"`
}

// interface for TextMessage and Message with attachment
type Message interface {
	Message()
}

// text message
type MessageText struct {
	Text string `json:"text"`
}

// message with attachment
type MessageAttachment struct {
	Attachment AttachmentFacebookRequest `json:"attachment"`
}

// implement Message
func (MessageText) Message()       {}
func (MessageAttachment) Message() {}

// attachment struct
type AttachmentFacebookRequest struct {
	AttachmentType string                    `json:"type"`
	Payload        AttachmentFacebookPayload `json:"payload"`
}

// payload struct
type AttachmentFacebookPayload struct {
	Src          string `json:"url,omitempty"`
	IsReusable   bool   `json:"is_reusable,omitempty"`
	TemplateType string `json:"template_type,omitempty"`
	Elements     []any  `json:"elements,omitempty"` // each element must match the template type
}

type Recipient struct {
	Id string `json:"id"`
}
