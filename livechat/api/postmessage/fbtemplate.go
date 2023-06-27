package postmessage

// A FBButton contains a facebook template button object informations.
type FBButton struct {
	URL   string `json:"url"`   // URL of the web_url type button action.
	Title string `json:"title"` // title of the button. (the text showed on button)
}

// A FBDefaultAction contains a facebook template's default action informations.
//
// Facebook template's default actions is an action perform when use click on any part of the template
// that doesn't have a button. Contains body close to that of a FBButton except for the lack of Title field.
type FBDefaultAction struct {
	URL string `json:"url"` // URL of the web_url type button action.
}

// A FBTemplateGeneric contains facebook generic template informations.
//
// *** request must have a title and one of the any other fields ***
//
// multiple buttons will align as a row
type FBTemplateGeneric struct {
	Title         string          `json:"title"`          // title of the template
	Message       string          `json:"message"`        // text message on the template
	Picture       string          `json:"picture"`        // image show on the template
	Button        []FBButton      `json:"buttons"`        // buttons shows on the template, maximum 3 buttons allow for a template
	DefaultAction FBDefaultAction `json:"default_action"` // the default action of the template
}

type FBTemplateButton struct{}

type FBTemplateCoupon struct{}

type FBTemplateCustomerFeedback struct{}

type FBTemplateMedia struct{}

type FBTemplateProduct struct{}

type FBTemplateReceipt struct{}

type FBTemplateStructuredInformation struct{}
