package reqigsendmessage

type DefaultAction struct {
	Type string `json:"type"` // web_url,postback
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
	Title         string        `json:"title"`
	ImageURL      string        `json:"image_url,omitempty"`
	Subtitle      string        `json:"subtitle,omitempty"`
	DefaultAction DefaultAction `json:"default_action,omitempty"`
	Buttons       []Button      `json:"buttons,omitempty"`
}

// TODO: define other template here

// implemet generic template
func (GenericTemplate) Template() {}
