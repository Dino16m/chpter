package apperrs

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ApplicationError struct {
	Code        codes.Code
	Description string
	Err         error
}

func (e ApplicationError) Error() string {
	return e.Description
}

func (e ApplicationError) Is(target error) bool {
	applicationError, ok := target.(ApplicationError)
	if !ok {
		return false
	}
	return applicationError.Code == e.Code || applicationError.Code == 0
}

func (e ApplicationError) GRPCStatus() *status.Status {
	return status.New(e.Code, e.Description)
}
