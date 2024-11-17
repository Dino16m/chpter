package apperrs

import (
	"errors"

	"google.golang.org/grpc/codes"
	"gorm.io/gorm"
)

var (
	ErrNotFound    ApplicationError = NewNotFoundError("not found")
	ErrServerError ApplicationError = NewServerError("internal server error")
)

func NewServerError(message string, errs ...error) ApplicationError {
	return createError(message, codes.Internal, errs...)
}

func NewNotFoundError(message string, errs ...error) ApplicationError {
	return createError(message, codes.NotFound, errs...)
}

func createError(message string, code codes.Code, errs ...error) ApplicationError {
	err := errors.Join(errs...)
	return ApplicationError{
		Code:        code,
		Description: message,
		Err:         err,
	}
}

func WrapNotFound[T any](result T, err error, message string) (T, error) {
	if err == nil {
		return result, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return result, NewNotFoundError(message, err)
	}
	return result, NewServerError("An error occurred", err)
}
