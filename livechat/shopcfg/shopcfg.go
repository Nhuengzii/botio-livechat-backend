// Package shopcfg defines the data model for a specific shop's config.
// It is to be stored in the database collection "shop_config".
package shopcfg

// Config is a shop's config data model
type Config struct {
	ShopID    string     `bson:"shopID"`
	Templates []Template `bson:"templates"`
}

// Template is a shop's saved template message data model
type Template struct {
	ID      string `bson:"id"`
	Payload string `bson:"payload"`
}
