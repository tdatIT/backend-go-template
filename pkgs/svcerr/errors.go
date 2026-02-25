package svcerr

import (
	"google.golang.org/grpc/codes"
)

// Error represents a business error with a specific error code and gRPC status
type Error struct {
	Message    string     // Error message
	VIMessage  string     // Optional message for VIM (if applicable)
	Code       string     // Business error code
	HTTPStatus int        // Optional HTTP status code for REST APIs
	GRPCCode   codes.Code // Optional gRPC status code for gRPC APIs
}

// Error implements the error interface
func (e *Error) Error() string {
	return e.Message
}
