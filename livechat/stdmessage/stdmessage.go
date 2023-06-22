package stdmessage

import (
	"fmt"
)

const (
	PlatformFacebook  Platform = "facebook"
	PlatformInstagram Platform = "instagram"
	PlatformLine      Platform = "line"
)

const (
	UserTypeUser  UserType = "user"
	UserTypeAdmin UserType = "admin"
)

const (
	AttachmentTypeImage                           AttachmentType = "image"
	AttachmentTypeVideo                           AttachmentType = "video"
	AttachmentTypeAudio                           AttachmentType = "audio"
	AttachmentTypeFile                            AttachmentType = "file"
	AttachmentTypeSticker                         AttachmentType = "sticker"
	AttachmentTypeLineEmoji                       AttachmentType = "line-emoji"
	AttachmentTypeLineTemplateButtons             AttachmentType = "line-template-buttons"
	AttachmentTypeLineTemplateConfirm             AttachmentType = "line-template-confirm"
	AttachmentTypeLineTemplateCarousel            AttachmentType = "line-template-carousel"
	AttachmentTypeLineTemplateImageCarousel       AttachmentType = "line-template-image-carousel"
	AttachmentTypeLineFlex                        AttachmentType = "line-flex"
	AttachmentTypeFBTemplateButton                AttachmentType = "facebook-template-button"
	AttachmentTypeFBTemplateCoupon                AttachmentType = "facebook-template-coupon"
	AttachmentTypeFBTemplateCustomerFeedback      AttachmentType = "facebook-template-customer-feedback"
	AttachmentTypeFBTemplateGeneric               AttachmentType = "facebook-template-generic"
	AttachmentTypeFBTemplateMedia                 AttachmentType = "facebook-template-media"
	AttachmentTypeFBTemplateProduct               AttachmentType = "facebook-template-product"
	AttachmentTypeFBTemplateReceipt               AttachmentType = "facebook-template-receipt"
	AttachmentTypeFBTemplateStructuredInformation AttachmentType = "facebook-template-structured-information"
	AttachmentTypeIGTemplateGeneric               AttachmentType = "instagram-template-generic"
	AttachmentTypeIGTemplateProduct               AttachmentType = "instagram-template-product"
)

type StdMessage struct {
	ShopID         string          `json:"shopID" bson:"shopID"`
	Platform       Platform        `json:"platform" bson:"platform"`
	PageID         string          `json:"pageID" bson:"pageID"`
	ConversationID string          `json:"conversationID" bson:"conversationID"`
	MessageID      string          `json:"messageID" bson:"messageID"`
	Timestamp      int64           `json:"timestamp" bson:"timestamp"`
	Source         Source          `json:"source" bson:"source"`
	Message        string          `json:"message" bson:"message"`
	Attachments    []Attachment    `json:"attachments" bson:"attachments"`
	ReplyTo        *RepliedMessage `json:"replyTo,omitempty" bson:"replyTo,omitempty"`
	IsDeleted      bool            `json:"isDeleted" bson:"isDeleted"`
}

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
		AttachmentTypeLineTemplateImageCarousel,
		AttachmentTypeFBTemplateButton,
		AttachmentTypeFBTemplateCoupon,
		AttachmentTypeFBTemplateCustomerFeedback,
		AttachmentTypeFBTemplateGeneric,
		AttachmentTypeFBTemplateMedia,
		AttachmentTypeFBTemplateProduct,
		AttachmentTypeFBTemplateReceipt,
		AttachmentTypeFBTemplateStructuredInformation,
		AttachmentTypeIGTemplateGeneric,
		AttachmentTypeIGTemplateProduct:
		return "ส่งเทมเพลต", nil
	case AttachmentTypeLineFlex:
		return "ส่งเฟล็กซ์", nil
	default:
		return "", fmt.Errorf("stdmessage.ToLastActivityString: unknown attachment type: %v", message.Attachments[0].AttachmentType)
	}
}
