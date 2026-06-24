package generatedownloadurl

import (
	"context"
	"time"

	"github.com/google/uuid"
	sharedv1 "github.com/zchelalo/neuraclinic-file-management/gen/go/shared/v1"
	fileerrors "github.com/zchelalo/neuraclinic-file-management/internal/modules/files/application/errors"
	appshared "github.com/zchelalo/neuraclinic-file-management/internal/modules/files/application/shared"
	"github.com/zchelalo/neuraclinic-file-management/internal/modules/files/ports"
)

type Command struct {
	ID uuid.UUID
}

type Result struct {
	DownloadURL string
	ExpiresAt   time.Time
}

type UseCase struct {
	cfg     appshared.Config
	repo    ports.Repository
	storage ports.Storage
}

func New(cfg appshared.Config, repo ports.Repository, storage ports.Storage) *UseCase {
	return &UseCase{cfg: cfg, repo: repo, storage: storage}
}

func (u *UseCase) Execute(ctx context.Context, cmd Command) (Result, error) {
	if cmd.ID == uuid.Nil {
		return Result{}, fileerrors.ErrInvalidInput
	}

	file, err := u.repo.ByID(ctx, cmd.ID)
	if err != nil {
		return Result{}, err
	}
	if file.Status != sharedv1.FileStatus_FILE_STATUS_AVAILABLE {
		return Result{}, fileerrors.ErrFailedPrecondition
	}

	downloadURL, expiresAt, err := u.storage.PresignDownload(ctx, file.StoragePath, u.cfg.DownloadURLTTL)
	if err != nil {
		return Result{}, err
	}

	return Result{DownloadURL: downloadURL, ExpiresAt: expiresAt}, nil
}
