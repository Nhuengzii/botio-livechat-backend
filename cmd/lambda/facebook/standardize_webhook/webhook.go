package main

type ReceiveWebhook struct {
	Object  string  `json:"object"`
	Entries []Entry `json:"entry"`
}

type Entry struct {
	PageID     string      `json:"id"`
	Time       int64       `json:"time"`
	Messagings []Messaging `json:"messaging"`
}

type Messaging struct {
	Sender    User    `json:"sender"`
	Recipient User    `json:"recipient"`
	Timestamp int64   `json:"timestamp"`
	Message   Message `json:"message"`
}

type Message struct {
	IsEcho      bool         `json:"is_echo"`
	MessageID   string       `json:"mid"`
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
	ReplyTo     ReplyMessage `json:"reply_to"`
}

type User struct {
	ID string `json:"id"`
}

type ReplyMessage struct {
	MessageId string `json:"mid"`
}

// ------ Attachment --------
type Attachment struct {
	AttachmentType string  `json:"type"`
	Payload        Payload `json:"payload"`
}
type Payload struct {
	Src string `json:"url,omitempty"` // for normal payload

	TemplateType string   `json:"template_type,omitempty"`
	Title        string   `json:"title,omitempty"`
	Subtitle     string   `json:"subtitle,omitempty"`
	ImageURL     string   `json:"image_url,omitempty"`
	Payload      string   `json:"payload,omitempty"`
	Buttons      []Button `json:"buttons,omitempty"` // for button template types
	Text         string   `json:"text,omitempty"`

	// --- Structured Information ---
	Countries []string `json:"countries"`
	//--- Coupon ---
	CouponCode           string `json:"coupon_code,omitempty"`
	CouponURL            string `json:"coupon_url,omitempty"`
	CouponURLButtonTitle string `json:"coupon_url_button_title,omitempty"`
	CouponPreMessage     string `json:"coupon_pre_message,omitempty"`

	// --- Customer Feedback ---
	ButtonTitle     string           `json:"button_title,omitempty"`
	FeedbackScreens []FeedbackScreen `json:"feedback_screens,omitempty"`
	BusinessPrivacy BusinessPrivacy  `json:"business_privacy,omitempty"`
	ExpiresInDays   string           `json:"expires_in_days,omitempty"`

	// --- Generic ---
	Elements []Element `json:"elements,omitempty"`

	// --- Receipt ---
	RecipientName string       `json:"recipient_name,omitempty"`
	OrderNumber   string       `json:"order_number,omitempty"`
	Currency      string       `json:"currency,omitempty"`
	PaymentMethod string       `json:"payment_method,omitempty"`
	OrderURL      string       `json:"order_url,omitempty"`
	Timestamp     int64        `json:"timestamp,omitempty"`
	Address       Address      `json:"address,omitempty"`
	Summary       Summary      `json:"summary,omitempty"`
	Adjustments   []Adjustment `json:"adjustments,omitempty"`
}

type Element struct {
	//--- Generic ---
	Title         string        `json:"title,omitempty"`
	Subtitle      string        `json:"subtitle,omitempty"`
	ImageURL      string        `json:"image_url,omitempty"`
	DefaultAction DefaultAction `json:"default_action,omitempty"`
	Buttons       []Button      `json:"buttons,omitempty"`

	// --- Reciept ---
	Quantity int    `json:"quantity,omitempty"`
	Price    int    `json:"price,omitempty"`
	Currency string `json:"currency,omitempty"`

	// --- Product ---
	ProductID string `json:"id,omitempty"`
	// --- Media ---
	MediaType    string `json:"media_type,omitempty"` // image,video
	AttachmentID string `json:"attachment_id,omitempty"`
}

type Button struct {
	Title string `json:"title"`
	Type  string `json:"type"` // web_url, postback, phone_number, account_link, account_unlink, game_play
	// web_url : url is the link
	// account_link : url is your login url
	Url                 string `json:"url,omitempty"`
	MessengerExtensions bool   `json:"messenger_extensions,omitempty"`
	WebviewHeightRatio  string `json:"webview_height_ratio,omitempty"`
	// postback : payload sent to postback webhook events
	// phone_number : payload is a phone number
	// game_play : payload is serialized json pauload
	Payload      string       `json:"payload,omitempty"`
	GameMetaData GameMetaData `json:"game_metadata,omitempty"`
}

type DefaultAction struct {
	Type                string `json:"type"`
	Url                 string `json:"url,omitempty"`
	MessengerExtensions bool   `json:"messenger_extensions,omitempty"`
	WebviewHeightRatio  string `json:"webview_height_ratio,omitempty"`
}

type FeedbackScreen struct {
	Questions []Question `json:"questions,omitempty"`
}

type Question struct {
	ID          string   `json:"id"`
	Type        string   `json:"type"`
	Title       string   `json:"title,omitempty"`
	ScoreLabel  string   `json:"score_label,omitempty"`
	ScoreOption string   `json:"score_option,omitempty"`
	FollowUp    FollowUp `json:"follow_up,omitempty"`
}

type FollowUp struct {
	Type        string `json:"type"`
	Placeholder string `json:"placeholder,omitempty"`
}
type BusinessPrivacy struct {
	URL string `json:"url"`
}

type Address struct {
	Street1    string `json:"street_1"`
	City       string `json:"city"`
	PostalCode string `json:"postal_code"`
	State      string `json:"state"`
	Country    string `json:"country"`
}

type Summary struct {
	SubTotal     string `json:"subtotal"`
	ShippingCost string `json:"shipping_cost"`
	TotalTax     string `json:"total_tax"`
	TotalCost    string `json:"total_cost"`
}
type Adjustment struct {
	Name   string `json:"name"`
	Amount int    `json:"amount"`
}
type GameMetaData struct {
	PlayerID  string `json:"player_id"`
	ContextID string `json:"context_id"`
}
