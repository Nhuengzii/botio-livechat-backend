package postmessage

// An InstagramButton contains a instagram template button object information.
type InstagramButton struct {
	URL   string `json:"url"`   // URL of the web_url type button action.
	Title string `json:"title"` // title of the button. (the text showed on button)
}

// An InstagramDefaultAction contains an instagram template's default action information.
//
// Instagram template's default actions is an action perform when use click on any part of the template
// that doesn't have a button. Contains body close to that of an InstagramButton except for the lack of Title field.
type InstagramDefaultAction struct {
	URL string `json:"url"` // URL of the web_url type button action.
}

// An InstagramTemplateGeneric contains instagram generic template information.
//
// *** request must have a title and at least one of other fields ***
//
// multiple buttons will align as a row
type InstagramTemplateGeneric struct {
	Title         string                  `json:"title"`                   // title of the template
	Message       string                  `json:"message,omitempty"`       // text message on the template
	Picture       string                  `json:"picture,omitempty"`       // image show on the template
	Button        []InstagramButton       `json:"buttons,omitempty"`       // buttons shows on the template, maximum 3 buttons allow for a template
	DefaultAction *InstagramDefaultAction `json:"defaultAction,omitempty"` // the default action of the template
}

type InstagramTemplateButton struct{}

type InstagramTemplateCoupon struct{}

type InstagramTemplateCustomerFeedback struct{}

type InstagramTemplateMedia struct{}

type InstagramTemplateProduct struct{}

type InstagramTemplateReceipt struct{}

type InstagramTemplateStructuredInformation struct{}
