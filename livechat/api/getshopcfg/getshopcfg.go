// Package getshopcfg defines the response data model for the API endpoint GET /shops/{shop_id}/config.
package getshopcfg

import "github.com/Nhuengzii/botio-livechat-backend/livechat/shopcfg"

// Response is a wrapper for shopcfg.Config used for the response of the API endpoint GET /shops/{shop_id}/config.
type Response struct {
	ShopConfig shopcfg.Config `json:"shopConfig"`
}
