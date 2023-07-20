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
	AttachmentTypeLineEmoji                       AttachmentType = "lineEmoji"
	AttachmentTypeLineTemplateButtons             AttachmentType = "lineTemplateButtons"
	AttachmentTypeLineTemplateConfirm             AttachmentType = "lineTemplateConfirm"
	AttachmentTypeLineTemplateCarousel            AttachmentType = "lineTemplateCarousel"
	AttachmentTypeLineTemplateImageCarousel       AttachmentType = "lineTemplateImageCarousel"
	AttachmentTypeLineFlex                        AttachmentType = "lineFlex"
	AttachmentTypeFBTemplateButton                AttachmentType = "facebookTemplateButton"
	AttachmentTypeFBTemplateCoupon                AttachmentType = "facebookTemplateCoupon"
	AttachmentTypeFBTemplateCustomerFeedback      AttachmentType = "facebookTemplateCustomerFeedback"
	AttachmentTypeFBTemplateGeneric               AttachmentType = "facebookTemplateGeneric"
	AttachmentTypeFBTemplateMedia                 AttachmentType = "facebookTemplateMedia"
	AttachmentTypeFBTemplateProduct               AttachmentType = "facebookTemplateProduct"
	AttachmentTypeFBTemplateReceipt               AttachmentType = "facebookTemplateReceipt"
	AttachmentTypeFBTemplateStructuredInformation AttachmentType = "facebookTemplateStructuredInformation"
	AttachmentTypeIGTemplateGeneric               AttachmentType = "instagramTemplateGeneric"
	AttachmentTypeIGTemplateProduct               AttachmentType = "instagramTemplateProduct"
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
//   - lineEmoji
//   - lineTemplateButtons
//   - lineTemplateConfirm
//   - lineTemplateCarousel
//   - lineTemplateImageCarousel
//   - lineFlex
//
// # Facebook specific AttachmentTypes.
//   - facebookTemplateButton
//   - facebookTemplateCoupon
//   - facebookTemplateCustomerFeedback.
//   - facebookTemplateGeneric
//   - facebookTemplateMedia
//   - facebookTemplateProduct
//   - facebookTemplateReceipt
//   - facebookTemplateStructuredInformation
//
// # Instagram specific AttachmentTypes.
//   - instagramTemplateGeneric
//   - instagramTemplateProduct
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
	lastActivity := ""
	if message.Source.UserType == UserType("admin") {
		lastActivity = "คุณ : "
	}

	if len(message.Attachments) == 0 {
		lastActivity += message.Message
		return lastActivity, nil
	}
	switch message.Attachments[0].AttachmentType {
	case AttachmentTypeImage:
		lastActivity += "ส่งรูปภาพ"
	case AttachmentTypeVideo:
		lastActivity += "ส่งวิดีโอ"
	case AttachmentTypeAudio:
		lastActivity += "ส่งข้อความเสียง"
	case AttachmentTypeFile:
		lastActivity += "ส่งไฟล์"
	case AttachmentTypeSticker:
		lastActivity += "ส่งสติกเกอร์"
	case AttachmentTypeLineEmoji:
		lastActivity += message.Message
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
		lastActivity += "ส่งเทมเพลต"
	case AttachmentTypeLineFlex:
		lastActivity += "ส่งเฟล็กซ์"
	default:
		return "", fmt.Errorf("stdmessage.ToLastActivityString: unknown attachment type: %v", message.Attachments[0].AttachmentType)
	}
	return lastActivity, nil
}
