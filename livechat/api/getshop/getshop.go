package getshop

import "github.com/Nhuengzii/botio-livechat-backend/livechat/shops"

type Response struct {
	AvailablePages []PlatformPageID `json:"availablePages"`
}

type PlatformPageID struct {
	PlatformName shops.Platform `json:"platformName"`
	PageID       string         `json:"pageID"`
}
