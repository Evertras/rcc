package main

import (
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"

	"github.com/evertras/rcc/pkg/repository"
	"github.com/evertras/rcc/pkg/server"
)

func main() {
	// TODO: proper config instead of magic env var
	tableName := os.Getenv("EVERTRAS_RCC_DYNAMODB_TABLE_NAME")
	if tableName == "" {
		panic("Missing table name")
	}

	repository, err := repository.NewDynamoDB(tableName)

	if err != nil {
		log.Fatal("Failed to create DynamoDB repository:", err)
	}

	server := server.New(server.NewDefaultConfig(), repository)

	// Even though we're behind a v2 gateway, we still use the v1 adapter here
	// as we seem to get the v1 event that contains the path, etc.
	lambda.Start(httpadapter.New(server.Handler()).ProxyWithContext)
}
