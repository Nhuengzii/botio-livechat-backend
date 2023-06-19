package livechat

import (
	"io"
)

type StorageUploader interface {
	UploadFile(bucketName string, file io.Reader) (string, error)
}
