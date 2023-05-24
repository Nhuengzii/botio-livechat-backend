package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/aws"
)

func sendMessage(svc *apigatewaymanagementapi.Client, connectionID string, message string) {
	input := &apigatewaymanagementapi.PostToConnectionInput{ConnectionId: aws.String(connectionID), Data: []byte(message)}
	_, err := svc.PostToConnection(context.Background(), input)
	if err != nil {
		discordLog(fmt.Sprint("Error sending message: ", err))
	}
}
