package user

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/balajiss36/Serverless/pkg/validators"
)

var (
	ErrorFailedToFetchRecord     = "failed to fetch record"
	ErrorInvalidUser             = "Invalid User"
	ErrorFailedToUnmarshalRecord = "Failed to UnMarshal Record"
	ErrorInvalidEmail            = "Email Invalid"
	ErrorCouldNotMarshalItem     = "Could not Marshal Item"
	ErrorCouldNotDeleteItem      = "Could not Delete Item"
	ErrorCouldNotDynamoPutItem   = "Could not put Dynamo Item"
	ErrorUserDoesNotExist        = "user does not exists"
	ErrorUserAlreadyExists       = "User Already Exists"
)

// Functions in handler will talk to the functions in the User.go to communicate with the database
// Each function will have 1:1 with the handler function

type User struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func FetchUser(email, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*User, error) { // Takes the value of email, tablename and dynaclient as per the handler.go
	input := &dynamodb.GetItemInput{ // Create a query to fetch the Item from the dynamodb based on the email
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		}, // From the db, match the key with the string which is an email
		// We go the pointer for Get Item Input on dynamodb, to get the string value for mem location of email we use aws.String
		TableName: aws.String(tableName),
	} // Same logic, input is given mem location for Get Item so to get the value from mem location of tableName we use aws.String
	result, err := dynaClient.GetItem(input) // Input is the query created for get the email and table name. using the dynaclient session created, we actually the value from dynamoDB
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	item := new(User)                                       // This is to UnMarshal the Go type into JSON so that the Struct User can understand and put into the struct User
	err = dynamodbattribute.UnmarshalMap(result.Item, item) // to Ummarshal the item to dynamodb https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/dynamodbattribute/
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}
	return item, nil

}

func FetchUsers(tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*[]User, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}
	result, err := dynaClient.Scan(input) // Scan is like GetAll, so we get get all users
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}
	item := new([]User)                                             // New returns the value as per the variable defined inside it
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, item) // to Ummarshal all the items of table from dynamodb so that it can be stored in a Go struct defined in item https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/dynamodbattribute/
	if err != nil {
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}
	return item, nil
}

func CreateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (
	*User,
	error,
) {
	var u User                                                   // We send the data we get from Client in User struct format
	if err := json.Unmarshal([]byte(req.Body), &u); err != nil { // From the request's body which has the JSON format, we pass it to u variable
		return nil, errors.New(ErrorInvalidUser)
	}
	if !validators.IsEmailValid(u.Email) { // Validate that the body of the request we put to 'u' var and validate that if it is a valid email address
		return nil, errors.New(ErrorInvalidEmail)
	}
	currentUser, _ := FetchUser(u.Email, tableName, dynaClient)
	if currentUser != nil && len(currentUser.Email) != 0 {
		return nil, errors.New(ErrorUserAlreadyExists)
	}
	av, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}
	_, err = dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}
	return &u, nil
}

func UpdateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (
	*User,
	error,
) {
	var u User
	if err := json.Unmarshal([]byte(req.Body), &u); err != nil {
		return nil, errors.New(ErrorInvalidEmail)
	} // Unmarshalled data is stored in the mem address of u
	// err := json.Unmarshal(cv, &us), cv is the []byte conversion to convert request's body into a slice of bytes for unmarshalling

	currentUser, _ := FetchUser(u.Email, tableName, dynaClient) // Run fetch User function to verify that the User exists or not.
	if currentUser != nil && len(u.Email) == 0 {
		return nil, errors.New(ErrorUserDoesNotExist)
	}

	av, err := dynamodbattribute.MarshalMap(u) // Marshal the user in a data which dynamodb understands so we use MarshalMap in the dynamodb attribute function
	if err != nil {
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}
	_, err = dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}
	return &u, nil
}

func DeleteUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) error {

	email := req.QueryStringParameters["email"]
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName),
	}
	_, err := dynaClient.DeleteItem(input)
	if err != nil {

		return errors.New(ErrorCouldNotDeleteItem)
	}
	return nil
}
