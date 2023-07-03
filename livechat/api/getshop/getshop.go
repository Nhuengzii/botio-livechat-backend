package getshop

import "github.com/Nhuengzii/botio-livechat-backend/livechat/shops"

type Response struct {
	AvailablePages []PlatformPageID `json:"available_pages"`
}

type PlatformPageID struct {
	PlatformName shops.Platform `json:"platform_name"`
	PageID       string         `json:"page_id"`
}
