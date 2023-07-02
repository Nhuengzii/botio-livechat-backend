package livechat

import "time"

// StorageUploader is an interface for storage uploader, for example, S3
type StorageClient interface {
	// UploadFile upload input file into the storage and return a location string if the operations was a success.
	// Return an error if it occurs.
	//
	// If implements by s3 storage name should be bucket name
	UploadFile(file []byte) (string, error)

	// RequestPutPresignedURL make a request to the storage and returns a PUT operation presigned URL.
	// Return an error if it occurs.
	//
	// PUT operation presignedURL can be used to upload a file to storage.
	// The URL is only valid for the time specified.
	RequestPutPresignedURL(validTime time.Duration) (string, error)
}
