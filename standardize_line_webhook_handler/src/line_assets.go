package main

import (
	"fmt"
	"reflect"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

// func getLineEmojiPNGurl(e *linebot.Emoji) string {
// 	return fmt.Sprintf("https://stickershop.line-scdn.net/sticonshop/v1/sticon/%s/android/%s.png", e.ProductID, e.EmojiID)
// }

func getStickerPNGurl(sm *linebot.StickerMessage) string {
	return fmt.Sprintf("https://stickershop.line-scdn.net/stickershop/v1/sticker/%s/android/sticker.png", sm.StickerID)
}

func getLocationString(lm *linebot.LocationMessage) string {
	return fmt.Sprintf("Title: %s\nAddress: %s\nLatitude: %f\nLongitude: %f", lm.Title, lm.Address, lm.Latitude, lm.Longitude)
}

// once the botio attachment structure supports line emojis embedding,
// then this function will be further worked on
func toLineEmojisBotioAttachments(tm *linebot.TextMessage) []attachment {
	return []attachment{}
}

func hasLineEmojis(tm *linebot.TextMessage) bool {
	v := reflect.ValueOf(tm).Elem().FieldByName("Emojis")
	return v != reflect.Value{}
}

func toStickerBotioAttachments(tm *linebot.StickerMessage) []attachment {
	var attachments = []attachment{}
	attachments = append(attachments,
		attachment{
			AttachmentType: attachmentTypeSticker,
			Payload: payload{
				Src: getStickerPNGurl(tm),
			},
		})
	return attachments
}
