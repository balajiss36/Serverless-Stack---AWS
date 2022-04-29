package handlers

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

// API request from User is marshaled from JSON to Go using stringBody.
func apiResponse(status int, body interface{}) (*events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{Headers: map[string]string{"Content-Type": "application/json"}}
	resp.StatusCode = status
	stringBody, _ := json.Marshal(body)
	resp.Body = string(stringBody)
	return &resp, nil
}
