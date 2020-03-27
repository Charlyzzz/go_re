package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
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

		resp, err := handler(events.APIGatewayProxyRequest{})
		if err != nil {
			t.Fatalf("Handler errored %v", err)
		}
		if statusCode := resp.StatusCode; statusCode != Redirect {
			t.Fatalf("Status Code should be %d but was %d", Redirect, statusCode)
		}
		if redirectionUri := resp.Headers[LocationHeader]; redirectionUri != google {
			t.Fatalf("Redirect URI should be %s but was %s", google, redirectionUri)
		}
	})

	t.Run("Not found", func(t *testing.T) {
		Finder = notFound()

		resp, err := handler(events.APIGatewayProxyRequest{})
		if err != nil {
			t.Fatalf("Handler errored %v", err)
		}
		if statusCode := resp.StatusCode; statusCode != NotFound {
			t.Fatalf("Status Code should be <%d> but was <%d>", NotFound, statusCode)
		}
		expectedBody := "not found"
		if body := resp.Body; body != expectedBody {
			t.Fatalf("Body should be <%s> but was <%s>", expectedBody, body)
		}
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
		if path := sr.Path; path != expectedPath {
			t.Fatalf("Path should be <%s> but was <%s>", expectedPath, path)
		}
		if subDomain := sr.SubDomain; subDomain != expectedSubDomain {
			t.Fatalf("SubDomain should be <%s> but was <%s>", expectedSubDomain, subDomain)
		}
	})
}
