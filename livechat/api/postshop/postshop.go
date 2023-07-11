// Package postshop defines the request and response data models for the API endpoint POST /shops.
package postshop

type Request struct {
	FacebookPageID       string `json:"facebookPageID"`
	FacebookAccessToken  string `json:"facebookAccessToken"`
	InstagramPageID      string `json:"instagramPageID"`
	InstagramAccessToken string `json:"instagramAccessToken"`
	LinePageID           string `json:"linePageID"`
	LineAccessToken      string `json:"lineAccessToken"`
	LineSecret           string `json:"lineSecret"`
}

// Response contains the shopID of the newly created shop.
// It is returned by the API endpoint POST /shops.
type Response struct {
	ShopID string `json:"shopID"`
}
