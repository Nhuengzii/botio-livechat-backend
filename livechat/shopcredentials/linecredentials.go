package shopcredentials

type LineCredentials struct {
	AccessToken string `bson:"accessToken"`
	Secret      string `bson:"secret"`
}
