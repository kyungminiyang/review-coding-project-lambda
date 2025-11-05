package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {
	tests := []struct {
		name           string
		request        events.LambdaFunctionURLRequest
		expectedStatus int
		expectedSize   int64
	}{
		{
			name: "Valid POST request with text file",
			request: events.LambdaFunctionURLRequest{
				RequestContext: events.LambdaFunctionURLRequestContext{
					HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
						Method: "POST",
					},
				},
				Body:            "Hello, this is test file content!",
				IsBase64Encoded: false,
				QueryStringParameters: map[string]string{
					"fileName": "test.txt",
				},
			},
			expectedStatus: 200,
			expectedSize:   33,
		},
		{
			name: "Valid POST request with base64 encoded file",
			request: events.LambdaFunctionURLRequest{
				RequestContext: events.LambdaFunctionURLRequestContext{
					HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
						Method: "POST",
					},
				},
				Body:            base64.StdEncoding.EncodeToString([]byte("Binary file content")),
				IsBase64Encoded: true,
				QueryStringParameters: map[string]string{
					"fileName": "test.bin",
				},
			},
			expectedStatus: 200,
			expectedSize:   19,
		},
		{
			name: "Valid POST request without fileName",
			request: events.LambdaFunctionURLRequest{
				RequestContext: events.LambdaFunctionURLRequestContext{
					HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
						Method: "POST",
					},
				},
				Body:            "Some content",
				IsBase64Encoded: false,
			},
			expectedStatus: 200,
			expectedSize:   12,
		},
		{
			name: "Invalid GET request",
			request: events.LambdaFunctionURLRequest{
				RequestContext: events.LambdaFunctionURLRequestContext{
					HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
						Method: "GET",
					},
				},
			},
			expectedStatus: 405,
		},
		{
			name: "Invalid PUT request",
			request: events.LambdaFunctionURLRequest{
				RequestContext: events.LambdaFunctionURLRequestContext{
					HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
						Method: "PUT",
					},
				},
			},
			expectedStatus: 405,
		},
		{
			name: "Empty file upload",
			request: events.LambdaFunctionURLRequest{
				RequestContext: events.LambdaFunctionURLRequestContext{
					HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
						Method: "POST",
					},
				},
				Body:            "",
				IsBase64Encoded: false,
				QueryStringParameters: map[string]string{
					"fileName": "empty.txt",
				},
			},
			expectedStatus: 200,
			expectedSize:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := Handler(context.Background(), tt.request)
			if err != nil {
				t.Fatalf("Handler returned an error: %v", err)
			}

			if response.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, response.StatusCode)
			}

			if tt.expectedStatus == 200 {
				var resp ReceiptResponse
				err := json.Unmarshal([]byte(response.Body), &resp)
				if err != nil {
					t.Fatalf("Failed to unmarshal response body: %v", err)
				}

				if resp.FileSize != tt.expectedSize {
					t.Errorf("Expected file size %d, got %d", tt.expectedSize, resp.FileSize)
				}

				if resp.Message == "" {
					t.Error("Expected a message in response")
				}

				if resp.Timestamp == 0 {
					t.Error("Expected a timestamp in response")
				}
			}
		})
	}
}
