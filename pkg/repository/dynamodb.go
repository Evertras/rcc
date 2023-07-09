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
		return nil, fmt.Errorf("session.NewSessionWithOptions: %w", err)
	}

	return &DynamoDB{
		tableName: tableName,
		dyndb:     dynamodb.New(sess),
	}, nil
}

func (d *DynamoDB) StoreValue1000(ctx context.Context, key string, value1000 int) error {
	if len(key) > 64 {
		return fmt.Errorf("key is too long")
	}

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
	if len(key) > 64 {
		return 0, fmt.Errorf("key is too long")
	}

	result, err := d.dyndb.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(d.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Key": {
				S: aws.String(key),
			},
		},
	})

	if err != nil {
		return 0, fmt.Errorf("d.dyndb.GetItem: %w", err)
	}

	if result.Item == nil {
		return 0, fmt.Errorf("not found")
	}

	entry := DynamoDBEntry{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &entry)
	if err != nil {
		return 0, fmt.Errorf("dynamodbattribute.UnmarshalMap: %w", err)
	}

	return entry.Value1000, nil
}
