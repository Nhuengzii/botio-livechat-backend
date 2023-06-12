package shops

type Shop struct {
	ShopID         string          `bson:"shopID"`
	FacebookPages  []FacebookPage  `bson:"facebookPages"`
	LinePages      []LinePage      `bson:"linePages"`
	InstagramPages []InstagramPage `bson:"instagramPages"`
}

type FacebookPage struct {
	PageID      string `bson:"pageID"`
	AccessToken string `bson:"accessToken"`
}

type LinePage struct {
	PageID      string `bson:"pageID"`
	AccessToken string `bson:"accessToken"`
	Secret      string `bson:"secret"`
}
type InstagramPage struct {
	PageID      string `bson:"pageID"`
	AccessToken string `bson:"accessToken"`
}
