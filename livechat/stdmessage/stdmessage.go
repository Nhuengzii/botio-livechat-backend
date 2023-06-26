// Package stdmessage define StdMessage a message format able to define various platform's messages.
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

// A StdMessage contains various information about specific message.
//
// # StdMessage is the main message's structure used to communicate in this system.
type StdMessage struct {
	ShopID         string          `json:"shopID" bson:"shopID"`                       // identified specific shop
	Platform       Platform        `json:"platform" bson:"platform"`                   // a platform the message belongs to
	PageID         string          `json:"pageID" bson:"pageID"`                       // identified specific page
	ConversationID string          `json:"conversationID" bson:"conversationID"`       // identified specific conversation
	MessageID      string          `json:"messageID" bson:"messageID"`                 // identified specific message
	Timestamp      int64           `json:"timestamp" bson:"timestamp"`                 // timestamp the message was sent
	Source         Source          `json:"source" bson:"source"`                       // author of the message
	Message        string          `json:"message" bson:"message"`                     // text component of the message
	Attachments    []Attachment    `json:"attachments" bson:"attachments"`             // various attachments sent with the message
	ReplyTo        *RepliedMessage `json:"replyTo,omitempty" bson:"replyTo,omitempty"` // this message was used to reply to other message
	IsDeleted      bool            `json:"isDeleted" bson:"isDeleted"`                 // if true the message was unsend and has been deleted
}

// A Platform used to define a platform that the message belongs to.
//
//   - facebook
//   - line
//   - instagram
type Platform string

// A Source store informations about the author of the message.
type Source struct {
	UserID string `json:"userID" bson:"userID"` // identified user
	// # user's type can be either user or admin.
	// 	- user : customer that has chat with the page.
	// 	- admin : page's admin.
	UserType UserType `json:"userType" bson:"userType"`
}

// # A UserType can be either
//   - user : customer that has chat with the page.
//   - admin : page's admin.
type UserType string

// An Attachment store various informations about specific attachment.
type Attachment struct {
	AttachmentType AttachmentType `json:"attachmentType" bson:"attachmentType"` // type of the attachment
	Payload        Payload        `json:"payload" bson:"payload"`               // actual payload of the attachment
}

// An AttachmentType stores attachment type string.
//
// # Non-specific-platform AttachmentTypes.
//   - image
//   - video
//   - audio
//   - file
//
// # Line specific AttachmentTypes.
//   - sticker
//   - line-emoji
//   - line-template-buttons
//   - line-template-confirm
//   - line-template-carousel
//   - line-template-image-carousel
//   - line-flex
//
// # Facebook specific AttachmentTypes.
//   - facebook-template-button
//   - facebook-template-coupon
//   - facebook-template-customer-feedback.
//   - facebook-template-generic
//   - facebook-template-media
//   - facebook-template-product
//   - facebook-template-receipt
//   - facebook-template-structured-information
//
// # Instagram specific AttachmentTypes.
//   - instagram-template-generic
//   - instagram-template-product
type AttachmentType string

// A Payload contains Src the content of media type attachments (image,video,audio,file).
type Payload struct {
	Src string `json:"src" bson:"src"` // URL storing payload's content
}

// A RepliedMessage contains information about target replied message.
type RepliedMessage struct {
	MessageID string `json:"messageID" bson:"messageID"` // identified message
}

// ToLastActivityString return a string that tells what activity last happen in the conversation.
// Return error if it occurs.
func (message *StdMessage) ToLastActivityString() (string, error) {
	if message.IsDeleted {
		return "ยกเลิกข้อความ", nil
	}
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
