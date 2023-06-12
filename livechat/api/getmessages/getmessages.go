package getmessages

import (
	"github.com/Nhuengzii/botio-livechat-backend/livechat/stdmessage"
)

type Response struct {
	Messages []*stdmessage.StdMessage `json:"messages"`
}
