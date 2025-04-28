package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

type Request struct {
	LicenseKey string `json:license_key`
	MachineID  string `json:machine_id`
}

type Response struct {
	Valid bool   `json:"valid"`
	Error string `json:"error,omitempty"`
}

func HandleRequest(ctx context.Context, req Request) (Response, error) {
	return Response{
		Valid: true,
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
