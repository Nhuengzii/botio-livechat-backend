package livechat

import (
	"io"
)

type StorageUploader interface {
	UploadFile(file io.Reader) (string, error)
}
