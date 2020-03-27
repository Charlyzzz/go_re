package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	. "go_re/common"
	"log"
	"net/http"
	"os"
	"time"
)

type CreateRecord struct {
	SubDomain   string `json:"subDomain" dynamodbav:"SubDomain"`
	Path        string `json:"path" dynamodbav:"Path"`
	RedirectUri string `json:"redirectUri" dynamodbav:"RedirectUri"`
}

var Saver RecordSaver

type saver struct {
	dynamoDB dynamodbiface.DynamoDBAPI
}

func NewRecordSaver() RecordSaver {
	config := &aws.Config{
		Region:   aws.String(os.Getenv("REGION_NAME")),
		Endpoint: aws.String(os.Getenv("DYNAMO_ENDPOINT")),
	}
	dynamoSession := session.Must(session.NewSession(config))
	return &saver{dynamoDB: dynamodb.New(dynamoSession)}
}

type RecordSaver interface {
	Call(createRecord CreateRecord) error
}

func (rc *saver) Call(createRecord CreateRecord) error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	item, err := dynamodbattribute.MarshalMap(createRecord)
	if err != nil {
		return err
	}
	log.Printf("%+v", item)

	putItemInput := &dynamodb.PutItemInput{
		TableName:           aws.String(os.Getenv("RECORDS_TABLE_NAME")),
		ConditionExpression: aws.String("attribute_not_exists(#subDomain)"),
		ExpressionAttributeNames: map[string]*string{
			"#subDomain": aws.String("SubDomain"),
		},
		Item: item,
	}
	_, err = rc.dynamoDB.PutItemWithContext(timeoutCtx, putItemInput)
	if err != nil {
		return err
	}
	return nil
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
	err = Saver.Call(*createRecord)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
	}, nil
}

func main() {
	Saver = NewRecordSaver()
	lambda.Start(handler)
}
