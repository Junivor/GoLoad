package errors

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewError(code codes.Code, format string, a ...interface{}) error {
	return status.Error(code, fmt.Sprintf(format, a...))
}

func ErrNotFound(entity string) error {
	return NewError(codes.NotFound, "%s not found", entity)
}

func ErrInternal(message string) error {
	return NewError(codes.Internal, message)
}
