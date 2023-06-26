package livechat

import (
	"io"
)

// StorageUploader is an interface for storage uploader, for example, S3
type StorageUploader interface {
	UploadFile(bucketName string, file io.Reader) (string, error)
}
