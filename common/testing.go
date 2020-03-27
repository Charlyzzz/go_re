package common

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Marshall(t *testing.T, any interface{}) string {
	assertions := assert.New(t)
	marshalledEntity, err := json.Marshal(any)
	assertions.NoError(err)
	return string(marshalledEntity)
}

func RunHandler(t *testing.T, handler Handler, request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	assertions := assert.New(t)
	resp, err := handler(request)
	assertions.NoError(err)
	return resp
}

func AssertHttpError(t *testing.T, response events.APIGatewayProxyResponse, statusCode int, httpError HttpError) {
	assertions := assert.New(t)
	assertions.Equal(statusCode, response.StatusCode)
	responseError := &HttpError{}
	err := json.Unmarshal([]byte(response.Body), responseError)
	assertions.NoError(err)
	assertions.EqualValues(httpError, *responseError)
}
