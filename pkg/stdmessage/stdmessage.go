package stdmessage

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

type Platform string

const (
	PlatformFacebook  Platform = "facebook"
	PlatformInstagram Platform = "instagram"
	PlatformLine      Platform = "line"
)

type Source struct {
	UserID   string   `json:"userID" bson:"userID"`
	UserType UserType `json:"userType" bson:"userType"`
}

type UserType string

const (
	UserTypeUser  UserType = "user"
	UserTypeGroup UserType = "group"
)

type Attachment struct {
	AttachmentType AttachmentType `json:"attachmentType" bson:"attachmentType"`
	Payload        Payload        `json:"payload" bson:"payload"`
}

type AttachmentType string

const (
	AttachmentTypeImage     AttachmentType = "image"
	AttachmentTypeVideo     AttachmentType = "video"
	AttachmentTypeAudio     AttachmentType = "audio"
	AttachmentTypeFile      AttachmentType = "file"
	AttachmentTypeSticker   AttachmentType = "sticker"
	AttachmentTypeLineEmoji AttachmentType = "line emoji"
)

type Payload struct {
	Src string `json:"src" bson:"src"`
}

type RepliedMessage struct {
	MessageID string `json:"messageID" bson:"messageID"`
}
