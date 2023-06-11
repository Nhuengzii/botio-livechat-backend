package getmessageresp

import "github.com/Nhuengzii/botio-livechat-backend/livechat"

type Response struct {
	Messages []*livechat.StdMessage `json:"messages"`
}
