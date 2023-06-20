package postmessage

type LineTemplateButtons struct {
	AltText           string               `json:"altText"`
	ThumbnailImageURL string               `json:"thumbnailImageUrl,omitempty"`
	Title             string               `json:"title,omitempty"`
	Text              string               `json:"text"`
	DefaultAction     LineTemplateAction   `json:"defaultAction,omitempty"`
	Actions           []LineTemplateAction `json:"actions"` // up to 4 actions
}

type LineTemplateConfirm struct {
	AltText string               `json:"altText"`
	Text    string               `json:"text"`
	Actions []LineTemplateAction `json:"actions"` // 2 actions only
}

type LineTemplateCarousel struct {
	AltText string                       `json:"altText"`
	Columns []LineTemplateCarouselColumn `json:"columns"` // up to 10 columns
}

type LineTemplateCarouselColumn struct {
	ThumbnailImageURL string               `json:"thumbnailImageUrl,omitempty"`
	Title             string               `json:"title,omitempty"`
	Text              string               `json:"text"`
	DefaultAction     LineTemplateAction   `json:"defaultAction,omitempty"`
	Actions           []LineTemplateAction `json:"actions"` // up to 3 actions
}

type LineTemplateImageCarousel struct {
	AltText string                            `json:"altText"`
	Columns []LineTemplateImageCarouselColumn `json:"columns"` // up to 10 columns
}

type LineTemplateImageCarouselColumn struct {
	ImageURL string             `json:"imageUrl"`
	Action   LineTemplateAction `json:"action"`
}

type LineTemplateAction struct {
	Label string `json:"label,omitempty"`
	URI   string `json:"uri"`
}
