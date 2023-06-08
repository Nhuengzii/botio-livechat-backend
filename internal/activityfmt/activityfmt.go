package activityfmt

import (
	"github.com/Nhuengzii/botio-livechat-backend/pkg/stdmessage"
)

func ToLastActivityString(message *stdmessage.StdMessage) string {
	if message.Message != "" {
		return message.Message
	}
	if len(message.Attachments) == 0 { // this really shouldn't be the case but just in case
		return ""
	}
	switch message.Attachments[0].AttachmentType {
	case stdmessage.AttachmentTypeImage:
		// return fmt.Sprintf("%s sent an image", displayName)
		return "new image message"
	case stdmessage.AttachmentTypeVideo:
		// return fmt.Sprintf("%s sent a video", displayName)
		return "new video message"
	case stdmessage.AttachmentTypeAudio:
		// return fmt.Sprintf("%s sent an audio", displayName)
		return "new audio message"
	case stdmessage.AttachmentTypeSticker:
		// return fmt.Sprintf("%s sent a sticker", displayName)
		return "new sticker message"
	default:
		return ""
	}
}
