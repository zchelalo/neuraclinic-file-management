package grpc

import (
	"errors"

	"github.com/zchelalo/neuraclinic-file-management/internal/modules/files/application"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapError(err error) error {
	switch {
	case errors.Is(err, application.ErrUnauthenticated):
		return status.Error(codes.Unauthenticated, "missing credentials")
	case errors.Is(err, application.ErrNotFound):
		return status.Error(codes.NotFound, "not found")
	case errors.Is(err, application.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, "invalid input")
	case errors.Is(err, application.ErrFailedPrecondition):
		return status.Error(codes.FailedPrecondition, "failed precondition")
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
