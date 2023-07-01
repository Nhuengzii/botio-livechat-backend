package reqigsendmessage

type DefaultAction struct {
	Type string `json:"type"` // web_url,postback URL  string `json:"url"`
	URL  string `json:"url"`
}
type Button struct {
	Type  string `json:"type"` // web_url
	URL   string `json:"url"`
	Title string `json:"title"`
}

type Template interface {
	Template()
}

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
