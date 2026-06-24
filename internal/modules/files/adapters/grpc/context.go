package grpc

import (
	"context"

	"github.com/google/uuid"
	fileerrors "github.com/zchelalo/neuraclinic-file-management/internal/modules/files/application/errors"
	"github.com/zchelalo/neuraclinic-file-management/internal/shared/appctx"
)

func userID(ctx context.Context) (uuid.UUID, error) {
	id, ok := appctx.UserID(ctx)
	if !ok || id == uuid.Nil {
		return uuid.Nil, fileerrors.ErrUnauthenticated
	}
	return id, nil
}

func parseID(value string) (uuid.UUID, error) {
	id, err := uuid.Parse(value)
	if err != nil || id == uuid.Nil {
		return uuid.Nil, fileerrors.ErrInvalidInput
	}
	return id, nil
}
