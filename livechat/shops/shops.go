package shops

type Shop struct {
	ShopID                  string                  `bson:"shopID"`
	FacebookPageID          string                  `bson:"facebookPageID"`
	LinePageID              string                  `bson:"linePageID"`
	InstagramPageID         string                  `bson:"instagramPageID"`
	FacebookAuthentication  FacebookAuthentication  `bson:"facebookAuthentication"`
	LineAuthentication      LineAuthentication      `bson:"lineAuthentication"`
	InstagramAuthentication InstagramAuthentication `bson:"instagramAuthentication"`
}

type FacebookAuthentication struct {
	AccessToken string `bson:"accessToken"`
}

type LineAuthentication struct {
	AccessToken string `bson:"accessToken"`
	Secret      string `bson:"secret"`
}
type InstagramAuthentication struct {
	AccessToken string `bson:"accessToken"`
}
