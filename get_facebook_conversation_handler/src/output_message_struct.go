package main

type OutputMessage struct {
	Conversations []Conversation `json:"conversations"`
}

type Conversation struct {
	ConversationID  string        `json:"conversationID"`
	ConversationPic Payload       `json:"conversationPic"`
	UpdatedTime     int64         `json:"updatedTime"`
	Participants    []Participant `json:"participants"`
	LastActivity    string        `json:"lastActivity"`
	IsRead          bool          `json:"isRead"`
}

type Participant struct {
	UserID     string  `json:"userID"`
	ProfilePic Payload `json:"profilePic"`
	Username   string  `json:"username"`
}
type Payload struct {
	Src string `json:"src"`
}
