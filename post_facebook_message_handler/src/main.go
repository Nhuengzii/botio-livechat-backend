package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler)
}

var (
	errNoPsidParam   = errors.New("QueryStringParameters psid not given")
	errNoPageIDParam = errors.New("PathParameter page_id not given")
)

func handler(context context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	start := time.Now()
	discordLog(fmt.Sprint("-------Post-FacebookMessage-handler!!!--------"))

	psid, ok := request.QueryStringParameters["psid"]
	if !ok {
		discordLog(fmt.Sprintf("Error reading psid queryStringParam"))
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, errNoPsidParam
	}
	pageID, ok := request.PathParameters["page_id"]
	if !ok {
		discordLog(fmt.Sprintf("Error reading pageID path param"))
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, errNoPageIDParam
	}

	var requestMessage RequestMessage
	err := json.Unmarshal([]byte(request.Body), &requestMessage)
	if err != nil {
		discordLog(fmt.Sprintf("Error unmarshal requestMessage : %v", err))
		return events.APIGatewayProxyResponse{}, err
	}

	var facebookResponse FacebookResponse
	err = SendFacebookMessage(requestMessage, psid, pageID, &facebookResponse)
	if err != nil {
		discordLog(fmt.Sprintf("Error sending facebook message : %v", err))
		return events.APIGatewayProxyResponse{}, err
	}
	discordLog(fmt.Sprintf("%+v", facebookResponse))
	discordLog(fmt.Sprintf("Elasped : %v", time.Since(start)))
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}, nil
}
