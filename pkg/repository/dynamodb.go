package repository

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DynamoDBEntry struct {
	Key       string
	Value1000 int
}

type DynamoDB struct {
	tableName string
}

func NewDynamoDB(tableName string) *DynamoDB {
	return &DynamoDB{
		tableName: tableName,
	}
}

func (d *DynamoDB) StoreValue1000(ctx context.Context, key string, value1000 int) error {
	// TODO: Do this once?
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})

	if err != nil {
		return fmt.Errorf("session.NewSessionWithOptions: %w", err)
	}

	svc := dynamodb.New(sess)

	item := DynamoDBEntry{
		Key:       key,
		Value1000: value1000,
	}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("dynamodbattribute.MarshalMap: %w", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(d.tableName),
	}

	_, err = svc.PutItem(input)

	if err != nil {
		return fmt.Errorf("svc.PutItem: %w", err)
	}

	return nil
}

func (d *DynamoDB) GetValue1000(ctx context.Context, key string) (int, error) {
	return 0, fmt.Errorf("not found")
}
