package shopcredentials

type FacebookCredentials struct {
	AccessToken         string `bson:"accessToken"`
	AppSecret           string `bson:"appSecret"`
	WebhookVerification string `bson:"webhookVerification"`
}
