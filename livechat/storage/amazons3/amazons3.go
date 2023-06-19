package amazons3

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
)

type Uploader struct {
	uploader *s3manager.Uploader
}

func NewUploader(awsRegion string) *Uploader {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	}))
	uploader := s3manager.NewUploader(sess)
	return &Uploader{uploader}
}

func (u *Uploader) UploadFile(bucketName string, file io.Reader) (string, error) {
	result, err := u.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Body:   file,
	})
	if err != nil {
		return "", fmt.Errorf("amazons3.Upload: %w", err)
	}
	return result.Location, nil
}
