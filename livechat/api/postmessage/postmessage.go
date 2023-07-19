// Package postmessage define request and response body for rest api endpoints that handle sending message request.
//
// The postmessage package should only be use to read api request body and create response body.
package postmessage

// A Request contains request body of the api request.
//
// *** only 1 field can exist at the same time. ***
type Request struct {
	Message    string     `json:"message"`    // send normal text message.
	Attachment Attachment `json:"attachment"` // send a message with an attachment.
}

// An Attachment contains attachment's details.
type Attachment struct {
	// should return error if the AttachmentType request does not match any of supported AttachmentType define in stdmessage package.
	AttachmentType string  `json:"type"`    // type of the attachment.
	Payload        Payload `json:"payload"` // store information of the attachment.
}

// A Payload store information of the attachment.
//
// *** only 1 field can exist at the same time. ***
//
// For facebook and instagram multiple templates specified will create a carousel.
//
// For line there is a template type for creating carousel.
type Payload struct {
	Src string `json:"src"` // contains URL if the AttachmentType are [image,audio,video,file].

	//[facebook templates]: https://developers.facebook.com/docs/messenger-platform/send-messages/templates
	//
	// Multiple templates specified will create a carousel
	FacebookTemplateGeneric []FacebookTemplateGeneric `json:"facebookTemplateGeneric"` // request's facebook generic template information.

	// [instagram generic template]: https://developers.facebook.com/docs/messenger-platform/instagram/features/generic-template
	// [instagram product template]: https://developers.facebook.com/docs/messenger-platform/instagram/features/product-template
	InstagramTemplateGeneric []InstagramTemplateGeneric `json:"instagramTemplateGeneric"` // request's instagram generic template information.

	// [line templates]: https://developers.line.biz/en/docs/messaging-api/message-types/#template-messages
	LineTemplateButtons       LineTemplateButtons       `json:"lineTemplateButtons"`       // request's line buttons template information.
	LineTemplateCarousel      LineTemplateCarousel      `json:"lineTemplateCarousel"`      // request's line carousel template information.
	LineTemplateImageCarousel LineTemplateImageCarousel `json:"lineTemplateImageCarousel"` // request's line image carousel template information.
	LineTemplateConfirm       LineTemplateConfirm       `json:"lineTemplateConfirm"`       // request's line confirm button template information.

	/* Templates below are currently unsupported
	FacebookTemplateButton                []FacebookTemplateButton                `json:"facebookTemplateButton"`
	FacebookTemplateCoupon                []FacebookTemplateCoupon                `json:"facebookTemplateCoupon"`
	FacebookTemplateCustomerFeedback      []FacebookTemplateCustomerFeedback      `json:"facebookTemplateCustomerFeedback"`
	FacebookTemplateProduct               []FacebookTemplateProduct               `json:"facebookTemplateProduct"`
	FacebookTemplateMedia                 []FacebookTemplateMedia                 `json:"facebookTemplateMedia"`
	FacebookTemplateReceipt               []FacebookTemplateReceipt               `json:"facebookTemplateReceipt"`
	FacebookTemplateStructuredInformation []FacebookTemplateStructuredInformation `json:"facebookTemplateStructuredInformation"`

	InstagramTemplateProduct []InstagramTemplateProduct `json:"instagramTemplateProduct"`
	*/
}

// A Response contains response body of the api request in case the call's a success.
type Response struct {
	RecipientID string `json:"recipientID"` // platform specific message receiver ID
	MessageID   string `json:"messageID"`   // platform specific requested message's message id
	Timestamp   int64  `json:"timestamp"`   // sent message's timestamp
}
