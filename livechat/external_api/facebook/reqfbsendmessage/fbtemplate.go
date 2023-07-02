package reqfbsendmessage

// A FBDefaultAction contains a facebook template's default action informations.
//
// Facebook template's default actions is an action perform when use click on any part of the template
// that doesn't have a button. Contains body close to that of a Button except for the lack of Title field.
//
// *** currently only support web_url type action ***
type DefaultAction struct {
	Type string `json:"type"` // web_url,postback
	URL  string `json:"url"`  // target URL if the Type of the button is web_url
}

// A Button contains a facebook template button object informations.
//
// *** currently only support web_url type action ***
type Button struct {
	Type  string `json:"type"`  // web_url,postback
	URL   string `json:"url"`   // target URL if the Type of the button is web_url
	Title string `json:"title"` // title of the button will be show as a text on the button
}

// Interface for various facebook templates
type Template interface {
	Template()
}

// A GenericTemplate contain informations of facebook Generic template type.“
//
// *** Required to have a Title and one of any other fields. ***
//
// [Facebook Generic Template doc] : https://developers.facebook.com/docs/messenger-platform/reference/templates/generic
type GenericTemplate struct {
	Title         string         `json:"title"`                    // title of the template will be show as big text on the template
	ImageURL      string         `json:"image_url,omitempty"`      // imageURL of the template's image
	Subtitle      string         `json:"subtitle,omitempty"`       // smaller text show on the template
	DefaultAction *DefaultAction `json:"default_action,omitempty"` // default action is an action that will be trigger if user click on any space the template except for buttons.
	Buttons       []Button       `json:"buttons,omitempty"`        // buttons shows on the template, maximum 3 buttons allow for a template
}

// TODO: define other template here

// implemet generic template
func (GenericTemplate) Template() {}
