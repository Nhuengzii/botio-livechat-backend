package stdconversation

type StdConversation struct {
	ShopID          string         `json:"shopID" bson:"shopID"`
	PageID          string         `json:"pageID" bson:"pageID"`
	ConversationID  string         `json:"conversationID" bson:"conversationID"`
	Platform        Platform       `json:"platform" bson:"platform"`
	ConversationPic Payload        `json:"conversationPic" bson:"conversationPic"`
	UpdatedTime     int64          `json:"updatedTime" bson:"updatedTime"`
	Participants    []*Participant `json:"participants" bson:"participants"`
	LastActivity    string         `json:"lastActivity" bson:"lastActivity"`
	IsRead          bool           `json:"isRead" bson:"isRead"`
}

type Participant struct {
	UserID     string  `json:"userID" bson:"userID"`
	Username   string  `json:"username" bson:"username"`
	ProfilePic Payload `json:"profilePic" bson:"profilePic"`
}

type Payload struct {
	Src string `json:"src" bson:"src"`
}

const (
	PlatformFacebook  Platform = "facebook"
	PlatformInstagram Platform = "instagram"
	PlatformLine      Platform = "line"
)

type Platform string