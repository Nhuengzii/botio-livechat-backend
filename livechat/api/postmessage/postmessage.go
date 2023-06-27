// Package postmessage define request and response body for rest api endpoints that handle sending message request.
//
// The postmessage package should only be use to read api request body and create response body.
package postmessage

// A Request contains request body of the api request.
//
// *** only 1 field can exists at the same time. ***
type Request struct {
	Message    string     `json:"message"`    // send normal text message.
	Attachment Attachment `json:"attachment"` // send a message with an attachment.
}

// An Attachment contains attachment's details.
type Attachment struct {
	// should return error if the AttachmentType request does not match any of supported AttachmentType define in stdmessage package.
	AttachmentType string  `json:"type"`    // type of the attachment.
	Payload        Payload `json:"payload"` // store informations of the attachment.
}

// A Payload store informations of the attachment.
//
// *** only 1 field can exists at the same time. ***
//
// For facebook and instagram multiple templates specified will create a carousel.
//
// For line there is a template type for creating carousel.
type Payload struct {
	Src string `json:"src"` // contains URL if the AttachmentType are [image,audio,video,file].

	//[facebook templates]: https://developers.facebook.com/docs/messenger-platform/send-messages/templates
	//
	// Multiple templates specified will create a carousel
	FBTemplateGeneric []FBTemplateGeneric `json:"fb_template_generic"` // request's facebook generic template informations.

	// [instagram generic template]: https://developers.facebook.com/docs/messenger-platform/instagram/features/generic-template
	// [instagram product template]: https://developers.facebook.com/docs/messenger-platform/instagram/features/product-template
	IGTemplateGeneric []IGTemplateGeneric `json:"ig_template_generic"` // request's intagram generic template informations.

	// [line templates]: https://developers.line.biz/en/docs/messaging-api/message-types/#template-messages
	LineTemplateButtons       LineTemplateButtons       `json:"line_template_buttons"`        // request's line buttons template informations.
	LineTemplateCarousel      LineTemplateCarousel      `json:"line_template_carousel"`       // request's line carousel template informations.
	LineTemplateImageCarousel LineTemplateImageCarousel `json:"line_template_image_carousel"` // request's line image carousel template informations.
	LineTemplateConfirm       LineTemplateConfirm       `json:"line_template_confirm"`        // request's line confirm button template informations.

	/* Templates below are currently unsupported
	FBTemplateButton                []FBTemplateButton                `json:"fb_template_button"`
	FBTemplateCoupon                []FBTemplateCoupon                `json:"fb_template_coupon"`
	FBTemplateCustomerFeedback      []FBTemplateCustomerFeedback      `json:"fb_template_customer_feedback"`
	FBTemplateMedia                 []FBTemplateMedia                 `json:"fb_template_media"`
	FBTemplateProduct               []FBTemplateProduct               `json:"fb_template_product"`
	FBTemplateReceipt               []FBTemplateReceipt               `json:"fb_template_receipt"`
	FBTemplateStructuredInformation []FBTemplateStructuredInformation `json:"fb_template_structured_information"`

	IGTemplateProduct []IGTemplateProduct `json:"ig_template_product"`
	*/
}

// A Response contains response body of the api request in case the call's a success.
type Response struct {
	RecipientID string `json:"recipient_id"` // platform specific message reciever ID
	MessageID   string `json:"message_id"`   // platform specific requested message's message id
	Timestamp   int64  `json:"timestamp"`    // sent message's timestamp
}
