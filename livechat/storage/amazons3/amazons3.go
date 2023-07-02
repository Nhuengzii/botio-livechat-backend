// Package amazons3 implements Uploader for manipulating amazons3 storage service
package amazons3

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
)

// An Uploader contains S3's session used to do various s3 bucket operations.
type Client struct {
	session    *session.Session
	bucketName string // target bucket name
}

// NewClient returns a new client which contains S3's session inside.
// Return an error if it occurs.
func NewClient(awsRegion string, bucketName string) *Client {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	}))
	return &Client{
		session:    sess,
		bucketName: bucketName,
	}
}

// UploadFile upload input file into the specified bucket and return a location string if the operations was a success.
// Return an error if it occurs.
func (c *Client) UploadFile(file []byte) (string, error) {
	log.Println(http.DetectContentType(file))

	uploader := s3manager.NewUploader(c.session)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(c.bucketName),
		Key:         aws.String(uuid.New().String()),
		ContentType: aws.String(http.DetectContentType(file)),
		Body:        bytes.NewReader(file),
	})
	if err != nil {
		return "", fmt.Errorf("amazons3.Upload: %w", err)
	}
	return result.Location, nil
}
