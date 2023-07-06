package livechat

import "time"

// StorageUploader is an interface for storage uploader, for example, S3
type StorageClient interface {
	// UploadFile upload input file into the storage and return a location string if the operations was a success.
	// Return an error if it occurs.
	//
	// If implements by s3 storage name should be bucket name
	UploadFile(file []byte) (string, error)

	// RequestPutPresignedURL make a request to S3 and returns a PUT operation presigned URL.
	// Returns URL for uploading to temporary storage s3 bucket if isTemporary is true.
	// Returns an error if it occurs.
	//
	// PUT operation presignedURL can be used to upload a file to Client's S3 bucket
	// The URL is only valid for the time specified.
	RequestPutPresignedURL(isTemporary bool, validDuration time.Duration) (string, error)
}
