package main

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	functions "github.com/aws/aws-sdk-go-v2/service/lambda"
	functionTypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context) (Response, error) {
	// Using the SDK's default configuration, loading additional config
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		// print error in golang
		panic(err)
	}

	// Using the Config value, create the DynamoDB client
	svc := functions.NewFromConfig(cfg)

	// invoke lambda function async
	_, errs := svc.Invoke(ctx, &functions.InvokeInput{
		FunctionName: aws.String("this-is-another-thing"),
		Payload:      []byte(`{"name": "John"}`),
		InvocationType: functionTypes.InvocationTypeEvent, // aws.String("Event"),
	})
	if errs != nil {
		panic(errs)
	}

	var buf bytes.Buffer

	body, err := json.Marshal(map[string]interface{}{
		"message": "Go Serverless v1.0! Your function executed successfully!",
	})
	if err != nil {
		return Response{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "hello-handler",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
