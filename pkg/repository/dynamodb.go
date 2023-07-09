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
	dyndb     *dynamodb.DynamoDB
}

func NewDynamoDB(tableName string) (*DynamoDB, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &DynamoDB{
		tableName: tableName,
		dyndb:     dynamodb.New(sess),
	}, nil
}

func (d *DynamoDB) StoreValue1000(ctx context.Context, key string, value1000 int) error {
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

	_, err = d.dyndb.PutItem(input)

	if err != nil {
		return fmt.Errorf("d.dyndb.PutItem: %w", err)
	}

	return nil
}

func (d *DynamoDB) GetValue1000(ctx context.Context, key string) (int, error) {
	return 0, fmt.Errorf("not found")
}
