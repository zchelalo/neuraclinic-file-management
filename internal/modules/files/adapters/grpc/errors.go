package grpc

import (
	"context"
	"errors"

	"github.com/zchelalo/neuraclinic-file-management/internal/modules/files/application"
	"github.com/zchelalo/neuraclinic-file-management/internal/shared/appctx"
	"github.com/zchelalo/neuraclinic-file-management/internal/shared/i18n"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapError(ctx context.Context, err error) error {
	language := appctx.Language(ctx)
	switch {
	case errors.Is(err, application.ErrUnauthenticated):
		return status.Error(codes.Unauthenticated, i18n.Message(language, i18n.KeyMissingCredentials))
	case errors.Is(err, application.ErrNotFound):
		return status.Error(codes.NotFound, i18n.Message(language, i18n.KeyNotFound))
	case errors.Is(err, application.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, i18n.Message(language, i18n.KeyInvalidInput))
	case errors.Is(err, application.ErrFailedPrecondition):
		return status.Error(codes.FailedPrecondition, i18n.Message(language, i18n.KeyFailedPrecondition))
	default:
		return status.Error(codes.Internal, i18n.Message(language, i18n.KeyInternalServerError))
	}
}
