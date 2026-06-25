package confirmupload

import (
	"context"
	"time"

	"github.com/google/uuid"
	sharedv1 "github.com/zchelalo/neuraclinic-file-management/gen/go/shared/v1"
	fileerrors "github.com/zchelalo/neuraclinic-file-management/internal/modules/files/application/errors"
	appshared "github.com/zchelalo/neuraclinic-file-management/internal/modules/files/application/shared"
	"github.com/zchelalo/neuraclinic-file-management/internal/modules/files/ports"
	"github.com/zchelalo/neuraclinic-file-management/internal/shared/appctx"
)

type Command struct {
	ID     uuid.UUID
	Status sharedv1.FileStatus
}

type Result struct {
	ID          uuid.UUID
	Status      sharedv1.FileStatus
	DownloadURL string
	ExpiresAt   time.Time
}

type UseCase struct {
	cfg       appshared.Config
	repo      ports.Repository
	storage   ports.Storage
	publisher ports.EventPublisher
	runtime   appshared.Runtime
}

func New(cfg appshared.Config, repo ports.Repository, storage ports.Storage, publisher ports.EventPublisher, runtime appshared.Runtime) *UseCase {
	return &UseCase{cfg: cfg, repo: repo, storage: storage, publisher: publisher, runtime: runtime.Normalize()}
}

func (u *UseCase) Execute(ctx context.Context, cmd Command) (Result, error) {
	if cmd.ID == uuid.Nil {
		return Result{}, fileerrors.ErrInvalidInput
	}
	if cmd.Status != sharedv1.FileStatus_FILE_STATUS_AVAILABLE && cmd.Status != sharedv1.FileStatus_FILE_STATUS_ERROR {
		return Result{}, fileerrors.ErrInvalidInput
	}

	file, err := u.repo.ByID(ctx, cmd.ID)
	if err != nil {
		return Result{}, err
	}
	if file.Status == sharedv1.FileStatus_FILE_STATUS_DELETED {
		return Result{}, fileerrors.ErrFailedPrecondition
	}

	if cmd.Status == sharedv1.FileStatus_FILE_STATUS_AVAILABLE {
		exists, err := u.storage.Exists(ctx, file.StoragePath)
		if err != nil {
			return Result{}, err
		}
		if !exists {
			return Result{}, fileerrors.ErrFailedPrecondition
		}
	}

	updated, err := u.repo.ConfirmUpload(ctx, cmd.ID, cmd.Status, u.runtime.Now().UTC())
	if err != nil {
		return Result{}, err
	}
	if err := u.publisher.PublishFileStatusChanged(ctx, ports.FileStatusChangedEvent{
		EventID:       u.runtime.NewUUID(),
		FileID:        updated.ID,
		ServiceOrigin: updated.ServiceOrigin,
		Status:        updated.Status,
		OccurredAt:    updated.UpdatedAt,
		RequestID:     appctx.RequestID(ctx),
		TraceID:       appctx.TraceID(ctx),
	}); err != nil {
		return Result{}, err
	}

	result := Result{
		ID:     updated.ID,
		Status: updated.Status,
	}
	if updated.Status == sharedv1.FileStatus_FILE_STATUS_AVAILABLE {
		result.DownloadURL, result.ExpiresAt, err = u.storage.PresignDownload(ctx, updated.StoragePath, u.cfg.DownloadURLTTL)
		if err != nil {
			return result, nil
		}
	}

	return result, nil
}
