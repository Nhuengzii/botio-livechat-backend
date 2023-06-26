// Package stdconversation define StdConversation a conversation format able to define various platform's conversation.
// It is stored and mainly used in this architecture.
package stdconversation

const (
	PlatformFacebook  Platform = "facebook"
	PlatformInstagram Platform = "instagram"
	PlatformLine      Platform = "line"
)

// A StdConversation contains various information about a specific conversation.
//
// # StdConversation is the main structure used to communicate in this system.
type StdConversation struct {
	ShopID          string        `json:"shopID" bson:"shopID"`                   // identified specific shop
	Platform        Platform      `json:"platform" bson:"platform"`               // specific platform the conversation belongs to
	PageID          string        `json:"pageID" bson:"pageID"`                   // identified specific page
	ConversationID  string        `json:"conversationID" bson:"conversationID"`   // identified specific conversation
	ConversationPic Payload       `json:"conversationPic" bson:"conversationPic"` // store conversation's display picture
	UpdatedTime     int64         `json:"updatedTime" bson:"updatedTime"`         // conversation's last message timestamp
	Participants    []Participant `json:"participants" bson:"participants"`       // an array of conversation's participants
	LastActivity    string        `json:"lastActivity" bson:"lastActivity"`       // last activity of the conversation
	Unread          int           `json:"unread" bson:"unread"`                   // number of conversation's currently unread messages
}

// A Platform used to define what platform the conversation belongs to.
type Platform string

// A Participant store informations about user participating in a specifc conversation
type Participant struct {
	UserID     string  `json:"userID" bson:"userID"`         // identified user
	Username   string  `json:"username" bson:"username"`     // user's username
	ProfilePic Payload `json:"profilePic" bson:"profilePic"` // user's profile pic store in payload object
}

// A Payload contains an attachment's URL
type Payload struct {
	Src string `json:"src" bson:"src"` // URL of the attachment
}
