// Package patchshop defines the request data model for the API endpoint PATCH /shops/{shop_id}.
package patchshop

// Request is the request body of the API endpoint PATCH /shops/{shop_id}.
type Request struct {
	FacebookPageID       string `json:"facebookPageID"`
	FacebookAccessToken  string `json:"facebookAccessToken"`
	InstagramPageID      string `json:"instagramPageID"`
	InstagramAccessToken string `json:"instagramAccessToken"`
	LinePageID           string `json:"linePageID"`
	LineAccessToken      string `json:"lineAccessToken"`
	LineSecret           string `json:"lineSecret"`
}
