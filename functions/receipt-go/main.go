package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// ReceiptResponse represents the API response structure
type ReceiptResponse struct {
	FileName  string `json:"fileName"`
	FileSize  int64  `json:"fileSize"`
	Timestamp int64  `json:"timestamp"`
	Message   string `json:"message"`
}

// Handler handles the Lambda function invocation
// Works with Lambda Function URLs
func Handler(ctx context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	// Only accept POST method
	if request.RequestContext.HTTP.Method != "POST" {
		return events.LambdaFunctionURLResponse{
			StatusCode: 405,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"Allow":        "POST",
			},
			Body: `{"error":"Method not allowed. Only POST is supported."}`,
		}, nil
	}

	// Check if body is base64 encoded
	var fileData []byte
	if request.IsBase64Encoded {
		decoded, err := base64.StdEncoding.DecodeString(request.Body)
		if err != nil {
			return events.LambdaFunctionURLResponse{
				StatusCode: 400,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				Body: `{"error":"Failed to decode base64 file data"}`,
			}, nil
		}
		fileData = decoded
	} else {
		fileData = []byte(request.Body)
	}

	// Get file size
	fileSize := int64(len(fileData))

	// Get file name from query parameters or headers
	fileName := "unknown"
	if request.QueryStringParameters != nil {
		if fn, ok := request.QueryStringParameters["fileName"]; ok && fn != "" {
			fileName = fn
		}
	}
	// Also check Content-Disposition header
	if request.Headers != nil {
		if cd, ok := request.Headers["content-disposition"]; ok && cd != "" {
			// Simple extraction, you can enhance this with regex
			fileName = cd
		}
	}

	// Create response
	receiptResponse := ReceiptResponse{
		FileName:  fileName,
		FileSize:  fileSize,
		Timestamp: time.Now().Unix(),
		Message:   "File received successfully",
	}

	// Marshal to JSON
	responseBytes, err := json.Marshal(receiptResponse)
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{"error":"Failed to generate response"}`,
		}, nil
	}

	return events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(responseBytes),
	}, nil
}

func main() {
	lambda.Start(Handler)
}

