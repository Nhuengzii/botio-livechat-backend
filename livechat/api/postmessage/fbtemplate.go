package postmessage

// A FacebookButton contains a facebook template button object information.
type FacebookButton struct {
	URL   string `json:"url"`   // URL of the web_url type button action.
	Title string `json:"title"` // title of the button. (the text showed on button)
}

// A FacebookDefaultAction contains a facebook template's default action information.
//
// Facebook template's default actions is an action perform when use click on any part of the template
// that doesn't have a button. Contains body close to that of a FacebookButton except for the lack of Title field.
type FacebookDefaultAction struct {
	URL string `json:"url"` // URL of the web_url type button action.
}

// A FacebookTemplateGeneric contains facebook generic template information.
//
// *** request must have a title and at least one of other fields ***
//
// multiple buttons will align as a row
type FacebookTemplateGeneric struct {
	Title         string                 `json:"title"`                   // title of the template
	Message       string                 `json:"message,omitempty"`       // text message on the template
	Picture       string                 `json:"picture,omitempty"`       // image show on the template
	Button        []FacebookButton       `json:"buttons,omitempty"`       // buttons shows on the template, maximum 3 buttons allow for a template
	DefaultAction *FacebookDefaultAction `json:"defaultAction,omitempty"` // the default action of the template
}

type FacebookTemplateButton struct{}

type FacebookTemplateCoupon struct{}

type FacebookTemplateCustomerFeedback struct{}

type FacebookTemplateMedia struct{}

type FacebookTemplateProduct struct{}

type FacebookTemplateReceipt struct{}

type FacebookTemplateStructuredInformation struct{}
