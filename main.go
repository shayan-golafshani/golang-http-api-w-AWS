package main

import (
	"fmt"
	awslambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/shayan-golafshani/golang-http-api-w-AWS/handlers"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/sendgrid/mcauto/apigw"
)

const (
	sgGatewayJWTKey        string = "jwt"
	sgGatewayResellerIDKey string = "resellerid"
	sgGatewayScopesKey     string = "scopes"
	sgGatewayUserIDKey     string = "userid"

	contentTypeHeader          = "Content-Type"
	contentTypeApplicationJSON = "application/json"
)

func main() {
	//load in credentials
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	server, err := handlers.Router()

	// Create DynamoDB client
	server.Db = dynamodb.New(sess)

	if err != nil {
		fmt.Println("Unable to create http router for lambda:httpServer", err)
	}
	awslambda.Start(apigw.Handle(server.Router))
}
