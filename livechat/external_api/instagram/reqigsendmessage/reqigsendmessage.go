package reqigsendmessage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func SendMessage(accessToken string, message SendingMessage, facebookPageID string) (_ *SendingMessageResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("reqigsendmessage.SendMessage: %w", err)
		}
	}()

	url := fmt.Sprintf("https://graph.facebook.com/v16.0/%v/messages?access_token=%v", facebookPageID, accessToken)
	instagramReqBody, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(instagramReqBody))
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

// instagram api response
type SendingMessageResponse struct {
	RecipientID string `json:"recipient_id"`
	MessageID   string `json:"message_id"`
	Timestamp   int64  `json:"timestamp"`
}

// instagram api send message top level struct
type SendingMessage struct {
	Recipient Recipient `json:"recipient"`
	Message   Message   `json:"message"`
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
	Attachment AttachmentInstagramRequest `json:"attachment"`
}

// implement Message
func (MessageText) Message()       {}
func (MessageAttachment) Message() {}

// attachment struct
type AttachmentInstagramRequest struct {
	AttachmentType string                     `json:"type"`
	Payload        AttachmentInstagramPayload `json:"payload"`
}

// payload struct
type AttachmentInstagramPayload struct {
	Src          string     `json:"url,omitempty"`
	IsReusable   bool       `json:"is_reusable,omitempty"`
	TemplateType string     `json:"template_type,omitempty"`
	Elements     []Template `json:"elements,omitempty"` // each element must match the template type
}

type Recipient struct {
	Id string `json:"id"`
}
