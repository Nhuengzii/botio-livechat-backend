// Package templates defines the data model for various templates.
// It is to be stored in the database collection "templates".
package templates

// Template is a shop's saved template message data model
type Template struct {
	ShopID  string `bson:"shopID" json:"shopID"`
	ID      string `bson:"templateID" json:"templateID"`
	Payload string `bson:"payload" json:"payload"`
}
