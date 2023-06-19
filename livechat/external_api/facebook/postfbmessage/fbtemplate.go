package postfbmessage

type DefaultAction struct {
	Type  string `json:"type"` // web_url,postback
	URL   string `json:"url"`
	Title string `json:"title"`
}
type Button struct {
	Type  string `json:"type"` // web_url
	URL   string `json:"url"`
	Title string `json:"title"`
}

type GenericTemplate struct {
	Title         string        `json:"title"`
	ImageURL      string        `json:"image_url"`
	Subtitle      string        `json:"subtitle"`
	DefaultAction DefaultAction `json:"default_action"`
	Buttons       []Button      `json:"buttons"`
}

// TODO: define other template here
