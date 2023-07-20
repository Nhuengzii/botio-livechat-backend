// Package stdconversation define StdConversation a conversation format able to define various platform's conversations.
package stdconversation

const (
	PlatformFacebook  Platform = "facebook"
	PlatformInstagram Platform = "instagram"
	PlatformLine      Platform = "line"
)

// A StdConversation contains various information about specific conversation.
//
// # StdConversation is the main conversation's structure used to communicate in this system.
type StdConversation struct {
	ShopID               string        `json:"shopID" bson:"shopID"`                             // identified specific shop
	Platform             Platform      `json:"platform" bson:"platform"`                         // a platform the conversation belongs to
	PageID               string        `json:"pageID" bson:"pageID"`                             // identified specific page
	ConversationID       string        `json:"conversationID" bson:"conversationID"`             // identified specific conversation
	ConversationPic      Payload       `json:"conversationPic" bson:"conversationPic"`           // store conversation's display picture
	UpdatedTime          int64         `json:"updatedTime" bson:"updatedTime"`                   // conversation's last message timestamp
	Participants         []Participant `json:"participants" bson:"participants"`                 // an array of conversation's participants
	LastActivity         string        `json:"lastActivity" bson:"lastActivity"`                 // last activity of the conversation
	LastUserActivityTime int64         `json:"lastUserActivityTime" bson:"lastUserActivityTime"` // last user activity time
	Unread               int           `json:"unread" bson:"unread"`                             // number of conversation's currently unread messages
}

// A Platform used to define a platform that the conversation belongs to.
//
//   - facebook
//   - line
//   - instagram
type Platform string

// A Participant store informations about user participating in a specifc conversation.
type Participant struct {
	UserID     string  `json:"userID" bson:"userID"`         // identified user
	Username   string  `json:"username" bson:"username"`     // user's username
	ProfilePic Payload `json:"profilePic" bson:"profilePic"` // user's profile pic store in payload object
}

// A Payload contains an attachment's URL
type Payload struct {
	Src string `json:"src" bson:"src"` // URL of the attachment
}
