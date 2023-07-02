package livechat

// StorageUploader is an interface for storage uploader, for example, S3
type StorageUploader interface {
	// UploadFile upload input file into the specified storage and return a location string if the operations was a success.
	// Return an error if it occurs.
	//
	// If implements by s3 storage name should be bucket name
	UploadFile(file []byte) (string, error)
}
