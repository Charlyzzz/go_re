package main

import (
	"github.com/stretchr/testify/assert"
	. "go_re/common"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

type fakeSaver struct {
	RecordSaver
}

func (fs *fakeSaver) Call(createRecord CreateRecord) error {
	return nil
}

func TestHandler(t *testing.T) {
	Saver = &fakeSaver{}

	t.Run("Without body", func(t *testing.T) {
		resp := RunHandler(t, handler, events.APIGatewayProxyRequest{})
		AssertHttpError(t, resp, http.StatusBadRequest, MissingBodyHttpError)
	})

	t.Run("Incomplete request", func(t *testing.T) {
		resp := RunHandler(t, handler, events.APIGatewayProxyRequest{
			Body: Marshall(t, &CreateRecord{}),
		})

		AssertHttpError(t, resp, http.StatusBadRequest, IncorrectBodyHttpError)
	})

	t.Run("Correct request", func(t *testing.T) {
		createRecord := &CreateRecord{
			SubDomain:   "google",
			Path:        "youtube",
			RedirectUri: "https://youtube.com",
		}

		resp := RunHandler(t, handler, events.APIGatewayProxyRequest{
			Body: Marshall(t, createRecord),
		})

		assertions := assert.New(t)
		assertions.Equal(http.StatusCreated, resp.StatusCode)
		assertions.Equal("", resp.Body)
	})
}
