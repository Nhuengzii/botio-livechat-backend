package postmessage

// An IGButton contains a instagram template button object informations.
type IGButton struct {
	URL   string `json:"url"`   // URL of the web_url type button action.
	Title string `json:"title"` // title of the button. (the text showed on button)
}

// An IGDefaultAction contains an instagram template's default action informations.
//
// Instagram template's default actions is an action perform when use click on any part of the template
// that doesn't have a button. Contains body close to that of an IGButton except for the lack of Title field.
type IGDefaultAction struct {
	URL string `json:"url"` // URL of the web_url type button action.
}

// An IGTemplateGeneric contains instagram generic template informations.
//
// *** request must have a title and one of the any other fields ***
//
// multiple buttons will align as a row
type IGTemplateGeneric struct {
	Title         string          `json:"title"`                    // title of the template
	Message       string          `json:"message,omitempty"`        // text message on the template
	Picture       string          `json:"picture,omitempty"`        // image show on the template
	Button        []IGButton      `json:"buttons,omitempty"`        // buttons shows on the template, maximum 3 buttons allow for a template
	DefaultAction IGDefaultAction `json:"default_action,omitempty"` // the default action of the template
}

type IGTemplateButton struct{}

type IGTemplateCoupon struct{}

type IGTemplateCustomerFeedback struct{}

type IGTemplateMedia struct{}

type IGTemplateProduct struct{}

type IGTemplateReceipt struct{}

type IGTemplateStructuredInformation struct{}
