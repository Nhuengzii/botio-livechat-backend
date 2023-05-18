package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var channelSecret = os.Getenv("CHANNEL_SECRET")
var qURL = os.Getenv("QUEUE_URL")

func signatureIsValid(channelSecret string, signature string, body []byte) bool {
	decoded, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false
	}
	hash := hmac.New(sha256.New, []byte(channelSecret))

	_, err = hash.Write(body)
	if err != nil {
		return false
	}

	return hmac.Equal(decoded, hash.Sum(nil))
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	body := []byte(req.Body)
	signature := req.Headers["x-line-signature"]
	if !signatureIsValid(channelSecret, signature, body) {
		log.Println("Unauthorized")
		return events.APIGatewayProxyResponse{
			StatusCode: 401,
			Body:       "Unauthorized",
		}, nil
	}
	sess := session.Must(session.NewSession())
	svc := sqs.New(sess, aws.NewConfig().WithRegion("ap-southeast-1"))
	params := &sqs.SendMessageInput{
		DelaySeconds: aws.Int64(0),
		MessageBody:  aws.String(string(body)),
		QueueUrl:     aws.String(qURL),
	}
	_, err := svc.SendMessage(params)
	if err != nil {
		log.Println("Cannot send message to SQS")
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Cannot send message to SQS",
		}, nil
	}
	log.Println("OK")
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "OK",
	}, nil
}

func main() {
	lambda.Start(Handler)
}
