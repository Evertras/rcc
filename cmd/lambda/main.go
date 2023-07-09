package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"

	"github.com/evertras/rcc/pkg/repository"
	"github.com/evertras/rcc/pkg/server"
)

func main() {
	repository := repository.NewInMemory()
	server := server.New(server.NewDefaultConfig(), repository)

	// Even though we're behind a v2 gateway, we still use the v1 adapter here
	// as we seem to get the v1 event that contains the path, etc.
	lambda.Start(httpadapter.New(server.Handler()).ProxyWithContext)
}
