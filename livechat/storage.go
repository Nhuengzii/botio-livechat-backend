package livechat

import (
	"io"
)

type StorageUploader interface {
	UploadFile(string, io.Reader) (string, error)
}
