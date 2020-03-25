package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const (
	Redirect        = 307
	NotFound        = 404
	LocationHeader  = "Location"
	SubDomainHeader = "SubDomainHeader"
)

var Finder RecordFinder

type SearchRecord struct {
	SubDomain string
	Path      string
}

func NewSearchRecord(request events.APIGatewayProxyRequest) SearchRecord {
	subDomain := request.Headers[SubDomainHeader]
	path := request.Path
	return SearchRecord{
		SubDomain: subDomain,
		Path:      path,
	}
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sr := NewSearchRecord(request)
	record, found, err := Finder.Call(sr)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	if !found {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       "not found",
		}, nil
	}
	return events.APIGatewayProxyResponse{
		StatusCode: Redirect,
		Headers:    map[string]string{LocationHeader: record.RedirectUri},
	}, nil
}

func main() {
	Finder = NewFinder()
	lambda.Start(handler)
}
