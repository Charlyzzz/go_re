package main

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func AssertString(t *testing.T, name string, expected string, given string) {
	if given != expected {
		t.Fatalf("%s should be <%s> but was <%s>", name, expected, given)
	}
}

func TestHandler(t *testing.T) {

	t.Run("Without body", func(t *testing.T) {
		resp, err := handler(events.APIGatewayProxyRequest{})
		if err != nil {
			t.Fatalf("Handler errored <%+v>", err)
		}
		if statusCode := resp.StatusCode; statusCode != http.StatusBadRequest {
			t.Fatalf("Status Code should be %d but was %d", http.StatusBadRequest, statusCode)
		}
		responseError := &HttpError{}
		err = json.Unmarshal([]byte(resp.Body), responseError)
		if err != nil {
			t.Fatalf("Unmarshalling error <%+v>", err)
		}
		missingBodyMessage := "Body is missing"
		if errorMessage := responseError.Error; errorMessage != missingBodyMessage {
			t.Fatalf("Body should be <%s> but was <%s>", missingBodyMessage, errorMessage)
		}
	})

	t.Run("Incomplete request", func(t *testing.T) {
		incompleteRequest := &CreateRecord{}
		incompleteRequestBody, err := json.Marshal(incompleteRequest)
		if err != nil {
			t.Fatalf("Marshalling error <%+v>", err)
		}
		resp, err := handler(events.APIGatewayProxyRequest{
			Body: string(incompleteRequestBody),
		})
		if err != nil {
			t.Fatalf("Handler errored <%+v>", err)
		}
		if statusCode := resp.StatusCode; statusCode != http.StatusBadRequest {
			t.Fatalf("Status Code should be %d but was %d", http.StatusBadRequest, statusCode)
		}
		responseError := &HttpError{}
		err = json.Unmarshal([]byte(resp.Body), responseError)
		if err != nil {
			t.Fatalf("Unmarshalling responseError <%+v>", err)
		}
		if errorMessage := responseError.Error; errorMessage != IncorrectBodyHttpError.Error {
			t.Fatalf("Body should be <%s> but was <%s>", IncorrectBodyHttpError.Error, errorMessage)
		}
	})

	t.Run("Correct request", func(t *testing.T) {
		createRecord := &CreateRecord{
			SubDomain:   "google",
			Path:        "youtube",
			RedirectUri: "https://youtube.com",
		}
		createRecordBody, err := json.Marshal(createRecord)
		if err != nil {
			t.Fatalf("Marshalling error <%+v>", err)
		}
		resp, err := handler(events.APIGatewayProxyRequest{
			Body: string(createRecordBody),
		})
		if err != nil {
			t.Fatalf("Handler errored <%+v>", err)
		}
		if statusCode := resp.StatusCode; statusCode != http.StatusCreated {
			t.Fatalf("Status Code should be %d but was %d", http.StatusCreated, statusCode)
		}
		if resp.Body != "" {
			t.Fatalf("Body should be empty")
		}
	})
}
