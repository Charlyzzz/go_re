package common

import "github.com/aws/aws-lambda-go/events"

type HttpError struct {
	Error string `json:"error"`
}

var (
	IncorrectBodyHttpError = HttpError{Error: "Body params missing or empty"}
	MissingBodyHttpError   = HttpError{Error: "Body is missing"}
)

type Handler = func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
