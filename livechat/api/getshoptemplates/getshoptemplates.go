// Package getshoptemplates defines the response data model for the API endpoint GET /shops/{shop_id}/config/templates.
package getshoptemplates

import "github.com/Nhuengzii/botio-livechat-backend/livechat/shopcfg"

type Response struct {
	Templates []shopcfg.Template `json:"templates"`
}
