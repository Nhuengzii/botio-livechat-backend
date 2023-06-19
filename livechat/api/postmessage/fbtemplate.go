package postmessage

type FBButton struct {
	URl   string `json:"url"`
	Title string `json:"title"`
}
type FBTemplateButton struct{}

type FBTemplateCoupon struct{}

type FBTemplateCustomerFeedback struct{}

type FBTemplateGeneric struct {
	Title         string     `json:"title"`
	Message       string     `json:"message"`
	Picture       string     `json:"picture"`
	Button        []FBButton `json:"buttons"`
	DefaultAction FBButton   `json:"default_action"`
}

type FBTemplateMedia struct{}

type FBTemplateProduct struct{}

type FBTemplateReceipt struct{}

type FBTemplateStructuredInformation struct{}
