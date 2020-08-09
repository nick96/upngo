package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const (
	signatureHeader = "X-Up-Authenticity-Signature"
	secretKeyEnvVar = "SECRET_KEY"
)

type WebhookEvent struct {
	Data struct {
		Type       string `json:"type"`
		ID         string `json:"id"`
		Attributes struct {
			EventType string `json:"eventType"`
			CreatedAt string `json:"createdAt"`
		} `json:"attributes"`
	}
}

func handleRequest(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Event: %v", event)
	if event.HTTPMethod == http.MethodGet {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusOK}, nil
	} else if event.HTTPMethod != http.MethodPost {
		log.Printf("Received request of method %s", event.HTTPMethod)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "POST only please",
		}, nil
	}

	expectedSignature, ok := event.Headers[signatureHeader]
	if !ok {
		log.Printf("No '%s' header in request", signatureHeader)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Nooooope",
		}, nil
	}

	secretKey := os.Getenv(secretKeyEnvVar)
	if strings.TrimSpace(secretKey) == "" {
		log.Printf("Secret key env var '%s' is empty or not provided", secretKeyEnvVar)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "We done goof",
		}, nil
	}

	mac := hmac.New(sha256.New, []byte(secretKey))
	_, err := io.Copy(mac, bytes.NewBufferString(event.Body))
	if err != nil {
		log.Printf("Failed to copy request body into mac")
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "We done goof",
		}, nil
	}
	actualSignature := mac.Sum(nil)
	if !hmac.Equal([]byte(expectedSignature), actualSignature) {
		log.Printf("Authenticity check of webhook request failed")
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "nice try buddy",
		}, nil
	}

	var webhookEvent WebhookEvent
	if err := json.Unmarshal([]byte(event.Body), &webhookEvent); err != nil {
		log.Printf("Request body isn't a valid WebhookEvent: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "up oh. That doesn't look right",
		}, nil
	}

	switch webhookEvent.Data.Type {
	case "PING":
		log.Printf("Received PING event")
		return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: "pong"}, nil
	default:
		log.Printf("Unhandled event: %s", webhookEvent.Data.Type)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError, Body: "Hmmm. Don't know about that"}, nil
	}
}

func main() {
	lambda.Start(handleRequest)
}
