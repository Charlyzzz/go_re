package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	. "go_re/common"
	"net/http"
	"testing"
)

type fakeFinder struct {
	record Record
	found  bool
	err    error
}

func found(record Record) RecordFinder {
	return &fakeFinder{
		record: record,
		found:  true,
	}
}

func notFound() RecordFinder {
	return &fakeFinder{}
}

func (fd *fakeFinder) Call(_ SearchRecord) (Record, bool, error) {
	if fd.err != nil {
		return Record{}, false, fd.err
	}
	var record Record
	if fd.found {
		record = fd.record
	} else {
		record = Record{}
	}
	return record, fd.found, nil
}

func TestHandler(t *testing.T) {

	t.Run("Redirection exists", func(t *testing.T) {
		google := "https://google.com"
		Finder = found(
			Record{
				SubDomain:   "NONE",
				Path:        "/search",
				RedirectUri: google,
			})

		resp := RunHandler(t, handler, events.APIGatewayProxyRequest{})

		assertions := assert.New(t)
		assertions.Equal(http.StatusTemporaryRedirect, resp.StatusCode)
		assertions.Equal(google, resp.Headers[LocationHeader])
	})

	t.Run("Not found", func(t *testing.T) {
		Finder = notFound()

		resp := RunHandler(t, handler, events.APIGatewayProxyRequest{})

		assertions := assert.New(t)
		assertions.Equal(http.StatusNotFound, resp.StatusCode)
		assertions.Equal("not found", resp.Body)
	})
}

func TestSearch(t *testing.T) {

	t.Run("New from request", func(t *testing.T) {

		expectedPath := "youtube"
		expectedSubDomain := "google"
		sr := NewSearchRecord(events.APIGatewayProxyRequest{
			PathParameters: map[string]string{
				SubDomain: "google",
				Path:      "youtube",
			},
		})

		assertions := assert.New(t)
		assertions.Equal(expectedPath, sr.Path)
		assertions.Equal(expectedSubDomain, sr.SubDomain)
	})
}
