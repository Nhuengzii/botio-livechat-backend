package amazons3

import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
)

type Uploader struct {
	uploader   *s3manager.Uploader
	bucketName string
}

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
