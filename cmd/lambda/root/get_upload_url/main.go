package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/api/getuploadurl"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/apigateway"
	"github.com/Nhuengzii/botio-livechat-backend/livechat/storage/amazons3"

	"github.com/Nhuengzii/botio-livechat-backend/livechat/discord"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const putPresignedURLValidDuration time.Duration = 10 * time.Minute

var errParsingTemporaryQueryParam = errors.New("temporary query string parameter must either be true or false")

func (c *config) handler(ctx context.Context, req events.APIGatewayProxyRequest) (_ events.APIGatewayProxyResponse, err error) {
	defer func() {
		if err != nil {
			logMessage := "cmd/lambda/root/get_upload_url/main.config.handler: " + err.Error()
			log.Println(logMessage)
			discord.Log(c.discordWebhookURL, logMessage)
		}
	}()

	queryStringParameters := req.QueryStringParameters
	isTemporary := false
	isTemporaryParamString, ok := queryStringParameters["temporary"]
	if ok {
		isTemporary, err = strconv.ParseBool(isTemporaryParamString)
		if err != nil {
			return apigateway.NewProxyResponse(400, fmt.Sprintf("Bad Request: %v", errParsingTemporaryQueryParam), "*"), nil
		}
	}

	presignedURL, err := c.storageClient.RequestPutPresignedURL(isTemporary, putPresignedURLValidDuration)
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
		s3TempBucketName  = os.Getenv("S3_TEMP_BUCKET_NAME")
	)
	storageClient := amazons3.NewClient(awsRegion, s3BucketName, s3TempBucketName)
	c := &config{
		discordWebhookURL: discordWebhookURL,
		awsRegion:         awsRegion,
		storageClient:     storageClient,
	}
	lambda.Start(c.handler)
}
