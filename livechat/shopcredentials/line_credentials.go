package shopcredentials

type LineCredentials struct {
	ChannelSecret      string `bson:"channelSecret"`
	ChannelAccessToken string `bson:"channelAccessToken"`
}
