// Package getshoptemplates defines the response data model for the API endpoint GET /shops/{shop_id}/config/templates.
package getshoptemplates

import (
	"github.com/Nhuengzii/botio-livechat-backend/livechat/templates"
)

type Response struct {
	Templates []templates.Template `json:"templates"`
}
