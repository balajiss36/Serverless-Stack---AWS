package main

import (
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/balajiss36/Serverless/pkg/handlers"
)

var (
	dynaClient dynamodbiface.DynamoDBAPI
)

func main() {
	region := os.Getenv("AWS_REGION") // We get the AWS region which is configured in the ENV variables while setting up AWS CLI, in my case us-west2
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region)})

	if err != nil {
		return
	}

	dynaClient = dynamodb.New(awsSession) // Start the dyanamodb service in the Session so that we can access the dynamodb base config
	lambda.Start(handler)                 // Start the lambda function as per the handler function
}

const tableName = "LambdaInGoUser"

func handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET": // Depending on the API method, we create different cases which will be executed as per the code in the handlers.go file
		return handlers.GetUser(req, tableName, dynaClient) // Get the table name and put the required to dyanamo service created.
	case "POST":
		return handlers.CreateUser(req, tableName, dynaClient) // Create the User  in the dyanmodb table
	case "PUT":
		return handlers.UpdateUser(req, tableName, dynaClient)
	case "DELETE":
		return handlers.DeleterUser(req, tableName, dynaClient)
	default:
		return handlers.UnhandledMethod()
	}

}
