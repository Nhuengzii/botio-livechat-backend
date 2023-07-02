package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/api/getuploadurl"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/apigateway"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/discord"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func (c *config) handler(ctx context.Context, req events.APIGatewayProxyRequest) (_ events.APIGatewayProxyResponse, err error) {
	defer func() {
		if err != nil {
			logMessage := "cmd/lambda/root/get_upload_url/main.config.handler: " + err.Error()
			log.Println(logMessage)
			discord.Log(c.discordWebhookURL, logMessage)
		}
	}()
	sess := session.Must(session.NewSession(aws.NewConfig().WithRegion(c.awsRegion)))
	svc := s3.New(sess)
	putObjReq, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(c.s3BucketName),
		Key:    aws.String(uuid.New().String()),
	})
	presignedURL, err := putObjReq.Presign(10 * time.Minute)
	if err != nil {
		return apigateway.NewProxyResponse(500, "Internal Server Error", "*"), err
	}
	response := getuploadurl.Response{
		PresignedURL: presignedURL,
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return apigateway.NewProxyResponse(500, "Internal Server Error", "*"), err
	}
	return apigateway.NewProxyResponse(200, string(responseJSON), "*"), nil
}

func main() {
	var (
		discordWebhookURL = os.Getenv("DISCORD_WEBHOOK_URL")
		awsRegion         = os.Getenv("AWS_REGION")
		s3BucketName      = os.Getenv("S3_BUCKET_NAME")
	)
	c := &config{
		discordWebhookURL: discordWebhookURL,
		awsRegion:         awsRegion,
		s3BucketName:      s3BucketName,
	}
	lambda.Start(c.handler)
}
