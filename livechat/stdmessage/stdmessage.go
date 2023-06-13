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

const (
	PlatformFacebook  Platform = "facebook"
	PlatformInstagram Platform = "instagram"
	PlatformLine      Platform = "line"
)

type Platform string

type Source struct {
	UserID   string   `json:"userID" bson:"userID"`
	UserType UserType `json:"userType" bson:"userType"`
}

const (
	UserTypeUser  UserType = "user"
	UserTypeAdmin UserType = "admin"
)

type UserType string

type Attachment struct {
	AttachmentType AttachmentType `json:"attachmentType" bson:"attachmentType"`
	Payload        Payload        `json:"payload" bson:"payload"`
}

const (
	AttachmentTypeImage                           AttachmentType = "image"
	AttachmentTypeVideo                           AttachmentType = "video"
	AttachmentTypeAudio                           AttachmentType = "audio"
	AttachmentTypeFile                            AttachmentType = "file"
	AttachmentTypeSticker                         AttachmentType = "sticker"
	AttachmentTypeLineEmoji                       AttachmentType = "lineEmoji"
	AttachmentTypeLineTemplateButtons             AttachmentType = "lineTemplateButtons"
	AttachmentTypeLineTemplateConfirm             AttachmentType = "lineTemplateConfirm"
	AttachmentTypeLineTemplateCarousel            AttachmentType = "lineTemplateCarousel"
	AttachmentTypeLineTemplateImageCarousel       AttachmentType = "lineTemplateImageCarousel"
	AttachmentTypeLineFlex                        AttachmentType = "lineFlex"
	AttachmentTypeFBTemplateButton                AttachmentType = "facebook-template-button"
	AttachmentTypeFBTemplateCoupon                AttachmentType = "facebook-template-coupon"
	AttachmentTypeFBTemplateCustomerFeedback      AttachmentType = "facebook-template-customerfeedback"
	AttachmentTypeFBTemplateGeneric               AttachmentType = "facebook-template-generic"
	AttachmentTypeFBTemplateMedia                 AttachmentType = "facebook-template-media"
	AttachmentTypeFBTemplateProduct               AttachmentType = "facebook-template-product"
	AttachmentTypeFBTemplateReceipt               AttachmentType = "facebook-template-receipt"
	AttachmentTypeFBTemplateStructuredInformation AttachmentType = "facebook-template-structuredinformation"
)

type AttachmentType string

type Payload struct {
	Src string `json:"src" bson:"src"`
}

type RepliedMessage struct {
	MessageID string `json:"messageID" bson:"messageID"`
}

var ErrUnknownAttachmentType = errors.New("stdmessage.ToLastActivityString: unknown attachment type")

func (message *StdMessage) ToLastActivityString() (string, error) {
	if len(message.Attachments) == 0 {
		return message.Message, nil
	}
	switch message.Attachments[0].AttachmentType {
	case AttachmentTypeImage:
		return "ส่งรูปภาพ", nil
	case AttachmentTypeVideo:
		return "ส่งวิดีโอ", nil
	case AttachmentTypeAudio:
		return "ส่งข้อความเสียง", nil
	case AttachmentTypeFile:
		return "ส่งไฟล์", nil
	case AttachmentTypeSticker:
		return "ส่งสติกเกอร์", nil
	case AttachmentTypeLineEmoji:
		return message.Message, nil
	case AttachmentTypeLineTemplateButtons,
		AttachmentTypeLineTemplateConfirm,
		AttachmentTypeLineTemplateCarousel,
		AttachmentTypeLineTemplateImageCarousel:
		return "ส่งเทมเพลท", nil
	case AttachmentTypeLineFlex:
		return "ส่งเฟล็กซ์", nil
	default:
		return "", ErrUnknownAttachmentType
	}
}
