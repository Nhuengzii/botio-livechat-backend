package main

import "github.com/aws/aws-lambda-go/lambda"

func main() {
	lambda.Start(handler)
}

func handler() {
	discordLog("get_facebook_conversation handler!!!")
}
