package main

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"log"
	"os"
	"time"
)

type Record struct {
	SubDomain   string
	Path        string
	RedirectUri string
}

func NewFinder() RecordFinder {
	config := &aws.Config{
		Region:   aws.String(os.Getenv("REGION_NAME")),
		Endpoint: aws.String(os.Getenv("DYNAMO_ENDPOINT")),
	}
	dynamoSession := session.Must(session.NewSession(config))
	return &finder{dynamoDB: dynamodb.New(dynamoSession)}
}

type finder struct {
	dynamoDB dynamodbiface.DynamoDBAPI
}

type RecordFinder interface {
	Call(searchRecord SearchRecord) (Record, bool, error)
}

func (rf *finder) Call(searchRecord SearchRecord) (Record, bool, error) {
	log.Printf("finding SubDomain<%s> Path<%s>", searchRecord.SubDomain, searchRecord.Path)
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	itemQuery := &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("RECORDS_TABLE_NAME")),
		Key: map[string]*dynamodb.AttributeValue{
			"Path": {
				S: aws.String(searchRecord.Path),
			},
			"SubDomain": {
				S: aws.String(searchRecord.SubDomain),
			},
		},
	}
	result, err := rf.dynamoDB.GetItemWithContext(timeoutCtx, itemQuery)
	record := Record{}
	if err != nil {
		return record, false, err
	}
	err = dynamodbattribute.UnmarshalMap(result.Item, &record)
	return record, record.RedirectUri != "", err
}
