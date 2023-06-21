package postmessage

type Request struct {
	Message    string     `json:"message"`
	Attachment Attachment `json:"attachment"`
}

type Attachment struct {
	AttachmentType string  `json:"type"`
	Payload        Payload `json:"payload"`
}

type Payload struct {
	Src string `json:"src"`

	FBTemplateButton                []FBTemplateButton                `json:"fb_template_button"`
	FBTemplateCoupon                []FBTemplateCoupon                `json:"fb_template_coupon"`
	FBTemplateCustomerFeedback      []FBTemplateCustomerFeedback      `json:"fb_template_customer_feedback"`
	FBTemplateGeneric               []FBTemplateGeneric               `json:"fb_template_generic"`
	FBTemplateMedia                 []FBTemplateMedia                 `json:"fb_template_media"`
	FBTemplateProduct               []FBTemplateProduct               `json:"fb_template_product"`
	FBTemplateReceipt               []FBTemplateReceipt               `json:"fb_template_receipt"`
	FBTemplateStructuredInformation []FBTemplateStructuredInformation `json:"fb_template_structured_information"`

	LineTemplateButtons       LineTemplateButtons       `json:"line_template_buttons"`
	LineTemplateConfirm       LineTemplateConfirm       `json:"line_template_confirm"`
	LineTemplateCarousel      LineTemplateCarousel      `json:"line_template_carousel"`
	LineTemplateImageCarousel LineTemplateImageCarousel `json:"line_template_image_carousel"`
}

type Response struct {
	RecipientID string `json:"recipient_id"`
	MessageID   string `json:"message_id"`
	Timestamp   int64  `json:"timestamp"`
}
