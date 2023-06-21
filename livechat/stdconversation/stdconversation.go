package stdconversation

const (
	PlatformFacebook  Platform = "facebook"
	PlatformInstagram Platform = "instagram"
	PlatformLine      Platform = "line"
)

type StdConversation struct {
	ShopID          string        `json:"shopID" bson:"shopID"`
	Platform        Platform      `json:"platform" bson:"platform"`
	PageID          string        `json:"pageID" bson:"pageID"`
	ConversationID  string        `json:"conversationID" bson:"conversationID"`
	ConversationPic Payload       `json:"conversationPic" bson:"conversationPic"`
	UpdatedTime     int64         `json:"updatedTime" bson:"updatedTime"`
	Participants    []Participant `json:"participants" bson:"participants"`
	LastActivity    string        `json:"lastActivity" bson:"lastActivity"`
	Unread          int           `json:"unread" bson:"unread"`
}

type Platform string

type Participant struct {
	UserID     string  `json:"userID" bson:"userID"`
	Username   string  `json:"username" bson:"username"`
	ProfilePic Payload `json:"profilePic" bson:"profilePic"`
}

type Payload struct {
	Src string `json:"src" bson:"src"`
}
