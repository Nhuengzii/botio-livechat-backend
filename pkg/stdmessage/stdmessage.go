package stdmessage

import "errors"

type StdMessage struct {
	ShopID         string          `json:"shopID" bson:"shopID"`
	Platform       Platform        `json:"platform" bson:"platform"`
	PageID         string          `json:"pageID" bson:"pageID"`
	ConversationID string          `json:"conversationID" bson:"conversationID"`
	MessageID      string          `json:"messageID" bson:"messageID"`
	Timestamp      int64           `json:"timestamp" bson:"timestamp"`
	Source         Source          `json:"source" bson:"source"`
	Message        string          `json:"message" bson:"message"`
	Attachments    []*Attachment   `json:"attachments,omitempty" bson:"attachments,omitempty"`
	ReplyTo        *RepliedMessage `json:"replyTo,omitempty" bson:"replyTo,omitempty"`
}

var (
	ErrNoAtttachment         = errors.New("Error ToLastActivityString no attachment in stdmessage")
	ErrUnknownAttachmentType = errors.New("Error ToLastActivityString unknown attachment type")
)

const (
	PlatformFacebook  Platform = "facebook"
	PlatformInstagram Platform = "instagram"
	PlatformLine      Platform = "line"
)

const (
	UserTypeUser  UserType = "user"
	UserTypeGroup UserType = "group"
)

const (
	AttachmentTypeImage     AttachmentType = "image"
	AttachmentTypeVideo     AttachmentType = "video"
	AttachmentTypeAudio     AttachmentType = "audio"
	AttachmentTypeFile      AttachmentType = "file"
	AttachmentTypeSticker   AttachmentType = "sticker"
	AttachmentTypeLineEmoji AttachmentType = "line emoji"
)

type Platform string

type Source struct {
	UserID   string   `json:"userID" bson:"userID"`
	UserType UserType `json:"userType" bson:"userType"`
}

type UserType string

type Attachment struct {
	AttachmentType AttachmentType `json:"attachmentType" bson:"attachmentType"`
	Payload        Payload        `json:"payload" bson:"payload"`
}

type AttachmentType string

type Payload struct {
	Src string `json:"src" bson:"src"`
}

type RepliedMessage struct {
	MessageID string `json:"messageID" bson:"messageID"`
}

func (message StdMessage) ToLastActivityString() (string, error) {
	if message.Message != "" {
		return message.Message, nil
	}
	if len(message.Attachments) == 0 { // this really shouldn't be the case but just in case
		return "", ErrNoAtttachment
	}
	switch message.Attachments[0].AttachmentType {
	case AttachmentTypeImage:
		// return fmt.Sprintf("%s sent an image", displayName)
		return "new image message", nil
	case AttachmentTypeVideo:
		// return fmt.Sprintf("%s sent a video", displayName)
		return "new video message", nil
	case AttachmentTypeAudio:
		// return fmt.Sprintf("%s sent an audio", displayName)
		return "new audio message", nil
	case AttachmentTypeSticker:
		// return fmt.Sprintf("%s sent a sticker", displayName)
		return "new sticker message", nil
	default:
		return "", ErrUnknownAttachmentType
	}
}
