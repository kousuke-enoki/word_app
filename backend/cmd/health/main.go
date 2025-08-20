// cmd/health/main.go
package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if req.Path == "/health" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       `{"ok":true}`,
		}, nil
	}
	return events.APIGatewayProxyResponse{StatusCode: 404, Body: "not found"}, nil
}

func main() { lambda.Start(handler) }
