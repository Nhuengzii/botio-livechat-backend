package main

type botioConversation struct {
	ShopID          string        `bson:"shopID"`
	PageID          string        `bson:"pageID"`
	ConversationID  string        `bson:"conversationID"`
	ConversationPic payload       `bson:"conversationPic"`
	UpdatedTime     int64         `bson:"updatedTime"`
	Participants    []participant `bson:"participants"`
	LastActivity    string        `bson:"lastActivity"`
	IsRead          bool          `bson:"isRead"` // this field is always false for LINE (unsupported)
}

type participant struct {
	UserID     string  `bson:"userID"`
	Username   string  `bson:"username"`
	ProfilePic payload `bson:"profilePic"`
}

func newBotioConversation(m *botioMessage) (*botioConversation, error) {
	userProfile, err := getLineUserProfile(m.Source.UserID)
	if err != nil {
		return nil, &newBotioConversationError{
			message: "couldn't create new botio conversation",
			err:     err,
		}
	}
	return &botioConversation{
		ShopID:          m.ShopID,
		PageID:          m.PageID,
		ConversationID:  m.ConversationID,
		ConversationPic: payload{Src: userProfile.PictureURL},
		UpdatedTime:     m.Timestamp,
		Participants: []participant{{
			UserID:     m.Source.UserID,
			Username:   userProfile.DisplayName,
			ProfilePic: payload{Src: userProfile.PictureURL},
		}},
		LastActivity: m.Message,
		IsRead:       false,
	}, nil
}

type newBotioConversationError struct {
	message string
	err     error
}

func (e *newBotioConversationError) Error() string {
	return e.message + ": " + e.err.Error()
}

func (e *newBotioConversationError) Unwrap() error {
	return e.err
}
