package utils

import (
	"errors"

	"github.com/sirupsen/logrus"
)

type Response struct {
	Status           string      `json:"status"`            // success or failed
	Message          string      `json:"message"`           // User-friendly message
	TechnicalMessage string      `json:"technical_message"` // Detailed technical info
	Data             interface{} `json:"data,omitempty"`    // Optional data
}

var responseStatus = []string{"success", "failed"}

// GenerateResponse creates a standardized response object
func GenerateResponse(status, message string, data interface{}, technicalMessage string) Response {
	// Validate status
	valid := false
	for _, s := range responseStatus {
		if s == status {
			valid = true
			break
		}
	}

	if !valid {
		logrus.Error("Invalid status provided.")
		panic(errors.New("Invalid status."))
	}

	// Validate message
	if message == "" {
		logrus.Error("Message is required.")
		panic(errors.New("Message is required."))
	}

	// Return a Response struct
	return Response{
		Status:           status,
		Message:          message,
		TechnicalMessage: technicalMessage,
		Data:             data,
	}
}
