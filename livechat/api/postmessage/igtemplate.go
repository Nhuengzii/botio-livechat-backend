package postmessage

type IGButton struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}
type IGDefaultAction struct {
	URL string `json:"url"`
}
type IGTemplateButton struct{}

type IGTemplateCoupon struct{}

type IGTemplateCustomerFeedback struct{}

type IGTemplateGeneric struct {
	Title         string          `json:"title"`
	Message       string          `json:"message"`
	Picture       string          `json:"picture"`
	Button        []FBButton      `json:"buttons"`
	DefaultAction FBDefaultAction `json:"default_action"`
}

type IGTemplateMedia struct{}

type IGTemplateProduct struct{}

type IGTemplateReceipt struct{}

type IGTemplateStructuredInformation struct{}
