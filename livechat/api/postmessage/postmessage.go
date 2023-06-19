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

	FBTemplateButton                []FBTemplateButton                `json:"fbTemplateButton"`
	FBTemplateCoupon                []FBTemplateCoupon                `json:"fbTemplateCoupon"`
	FBTemplateCustomerFeedback      []FBTemplateCustomerFeedback      `json:"fbTemplateCustomerFeedback"`
	FBTemplateGeneric               []FBTemplateGeneric               `json:"fbTemplateGeneric"`
	FBTemplateMedia                 []FBTemplateMedia                 `json:"fbTemplateMedia"`
	FBTemplateProduct               []FBTemplateProduct               `json:"fbTemplateProduct"`
	FBTemplateReceipt               []FBTemplateReceipt               `json:"fbTemplateReceipt"`
	FBTemplateStructuredInformation []FBTemplateStructuredInformation `json:"fbTemplateStructuredInformation"`
}

type Response struct {
	RecipientID string `json:"recipient_id"`
	MessageID   string `json:"message_id"`
	Timestamp   int64  `json:"timestamp"`
}
