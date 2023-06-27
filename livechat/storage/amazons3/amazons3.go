// Package amazons3 implements Uploader for manipulating amazons3 storage service
package amazons3

import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
)

// An Uploader contains s3manager.Uploader used to do various s3 bucket operations.
type Uploader struct {
	uploader   *s3manager.Uploader // s3 sdk
	bucketName string              // target bucket name
}

// NewUploader returns a new Uploader which contains s3 uploader inside.
// Return an error if it occurs.
func NewUploader(awsRegion string, bucketName string) *Uploader {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	}))
	uploader := s3manager.NewUploader(sess)
	return &Uploader{
		uploader:   uploader,
		bucketName: bucketName,
	}
}

// UploadFile upload input file into the specified bucket and return a location string if the operations was a success.
// Return an error if it occurs.
func (u *Uploader) UploadFile(file io.Reader) (string, error) {
	result, err := u.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(u.bucketName),
		Key:    aws.String(uuid.New().String()),
		Body:   file,
	})
	if err != nil {
		return "", fmt.Errorf("amazons3.Upload: %w", err)
	}
	return result.Location, nil
}
