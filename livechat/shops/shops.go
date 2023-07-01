// Package shops defines Shop, a struct storing config about a specific shop
package shops

type Shop struct {
	ShopID                  string                  `bson:"shopID"`                  // identified specific shop
	FacebookPageID          string                  `bson:"facebookPageID"`          // identified specific facebook page
	LinePageID              string                  `bson:"linePageID"`              // identified specific line page
	InstagramPageID         string                  `bson:"instagramPageID"`         // identified specific instagram page
	FacebookAuthentication  FacebookAuthentication  `bson:"facebookAuthentication"`  // store page's facebook authentication variables
	LineAuthentication      LineAuthentication      `bson:"lineAuthentication"`      // store page's line authentication variables
	InstagramAuthentication InstagramAuthentication `bson:"instagramAuthentication"` // store page's instagram authentication variables
}

// A FacebookAuthentication store page's facebook authentication variables
type FacebookAuthentication struct {
	AccessToken string `bson:"accessToken"` // page's access token
}

// A LineAuthentication store page's line authentication variables
type LineAuthentication struct {
	AccessToken string `bson:"accessToken"` // line's bot access token
	Secret      string `bson:"secret"`      // line's app secret
}

// An InstagramAuthentication store page's instagram authentication variables
type InstagramAuthentication struct {
	AccessToken string `bson:"accessToken"` // instagram's page access token
}

type Platform string

const (
	PlatformFacebook  Platform = "facebook"
	PlatformInstagram Platform = "instagram"
	PlatformLine      Platform = "line"
)
