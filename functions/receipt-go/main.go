package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"vibe-coding-project-lambda/shared/repository"
)

const (
	defaultBucketName = "review-receipt-uploads"
	defaultRegion     = "ap-northeast-1"
)

var s3Repo *repository.S3Repository

// init initializes S3 repository
func init() {
	ctx := context.Background()

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(defaultRegion))
	if err != nil {
		panic(fmt.Sprintf("unable to load SDK config: %v", err))
	}

	// Initialize S3 client
	s3Client := s3.NewFromConfig(cfg)

	// Get bucket name from environment variable or use default
	bucketName := os.Getenv("S3_BUCKET_NAME")
	if bucketName == "" {
		bucketName = defaultBucketName
	}

	// Create repository layer
	s3Repo = repository.NewS3Repository(s3Client, bucketName, defaultRegion)
}

// UploadRequest represents the file upload request structure
type UploadRequest struct {
	FileName    string `json:"filename"`
	FileContent string `json:"file_content"` // Base64 encoded file content
	ContentType string `json:"content_type"`
}

// ReceiptResponse represents the API response structure
type ReceiptResponse struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	FileName   string `json:"fileName,omitempty"`
	FileSize   int64  `json:"fileSize,omitempty"`
	BucketName string `json:"bucketName,omitempty"`
	Key        string `json:"key,omitempty"`
	URL        string `json:"url,omitempty"`
	UploadDate string `json:"uploadDate,omitempty"`
	Timestamp  int64  `json:"timestamp"`
	Error      string `json:"error,omitempty"`
}

// Handler handles the Lambda function invocation
// Works with Lambda Function URLs
func Handler(ctx context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	timestamp := time.Now().Unix()

	// Handle CORS preflight requests
	if request.RequestContext.HTTP.Method == "OPTIONS" {
		return events.LambdaFunctionURLResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "POST, OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type",
			},
		}, nil
	}

	// Only accept POST method
	if request.RequestContext.HTTP.Method != "POST" {
		return errorResponse(405, "Method not allowed. Only POST is supported.", "", timestamp)
	}

	// Parse JSON request body
	var uploadReq UploadRequest
	if err := json.Unmarshal([]byte(request.Body), &uploadReq); err != nil {
		return errorResponse(400, "Failed to parse request body", err.Error(), timestamp)
	}

	// Validate required fields
	if uploadReq.FileName == "" || uploadReq.FileContent == "" {
		return errorResponse(400, "Missing required fields: filename and file_content", "", timestamp)
	}

	// Set default content type if not provided
	if uploadReq.ContentType == "" {
		uploadReq.ContentType = "application/octet-stream"
	}

	// Decode base64 file content
	fileData, err := base64.StdEncoding.DecodeString(uploadReq.FileContent)
	if err != nil {
		return errorResponse(400, "Failed to decode base64 file data", err.Error(), timestamp)
	}

	// Validate we have file content
	if len(fileData) == 0 {
		return errorResponse(400, "File content is empty", "", timestamp)
	}

	// Upload to S3
	fileInfo, err := s3Repo.Upload(ctx, uploadReq.FileName, fileData, uploadReq.ContentType)
	if err != nil {
		log.Printf("Failed to upload to S3: %v", err)
		return errorResponse(500, "Failed to upload file to S3", err.Error(), timestamp)
	}

	// Create success response
	response := ReceiptResponse{
		Success:    true,
		Message:    "File uploaded successfully to S3",
		FileName:   fileInfo.FileName,
		FileSize:   fileInfo.Size,
		BucketName: fileInfo.BucketName,
		Key:        fileInfo.Key,
		URL:        fileInfo.URL,
		UploadDate: fileInfo.UploadDate,
		Timestamp:  timestamp,
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		return errorResponse(500, "Failed to generate response", err.Error(), timestamp)
	}

	return events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "POST, OPTIONS",
			"Access-Control-Allow-Headers": "Content-Type",
		},
		Body: string(responseBytes),
	}, nil
}

// errorResponse creates an error response
func errorResponse(statusCode int, message, errorDetail string, timestamp int64) (events.LambdaFunctionURLResponse, error) {
	response := ReceiptResponse{
		Success:   false,
		Message:   message,
		Error:     errorDetail,
		Timestamp: timestamp,
	}

	responseBody, _ := json.Marshal(response)

	headers := map[string]string{
		"Content-Type":                 "application/json",
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "POST, OPTIONS",
		"Access-Control-Allow-Headers": "Content-Type",
	}
	if statusCode == 405 {
		headers["Allow"] = "POST"
	}

	return events.LambdaFunctionURLResponse{
		StatusCode: statusCode,
		Headers:    headers,
		Body:       string(responseBody),
	}, nil
}

func main() {
	lambda.Start(Handler)
}
