package postmessage

type LineMessageType string

const (
	LineMessageTypeTemplate LineMessageType = "template"
)

type LineTemplateMessage struct {
	Type     LineMessageType `json:"type"`
	AltText  string          `json:"altText"`
	Template interface{}     `json:"template"`
}
