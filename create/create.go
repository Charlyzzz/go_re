package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	. "go_re/common"
	"net/http"
)

type CreateRecord struct {
	SubDomain   string `json:"subDomain"`
	Path        string `json:"path"`
	RedirectUri string `json:"redirectUri"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.Body == "" {
		errorResponse, err := json.Marshal(MissingBodyHttpError)
		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       string(errorResponse),
		}, nil
	}
	createRecord := &CreateRecord{}
	err := json.Unmarshal([]byte(request.Body), createRecord)
	if createRecord.Path == "" || createRecord.SubDomain == "" || createRecord.RedirectUri == "" {
		errorResponse, err := json.Marshal(&IncorrectBodyHttpError)
		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       string(errorResponse),
		}, nil
	}
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
	}, nil
}

func main() {
	lambda.Start(handler)
}
