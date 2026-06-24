package requestupload

import (
	"context"
	"path"
	"strings"
	"time"

	"github.com/google/uuid"
	sharedv1 "github.com/zchelalo/neuraclinic-file-management/gen/go/shared/v1"
	fileerrors "github.com/zchelalo/neuraclinic-file-management/internal/modules/files/application/errors"
	appshared "github.com/zchelalo/neuraclinic-file-management/internal/modules/files/application/shared"
	"github.com/zchelalo/neuraclinic-file-management/internal/modules/files/ports"
)

type Command struct {
	UploadedBy    uuid.UUID
	OriginalName  string
	MimeType      string
	SizeBytes     int64
	IsPublic      bool
	ServiceOrigin string
}

type Result struct {
	ID        uuid.UUID
	UploadURL string
	ExpiresAt time.Time
}

type UseCase struct {
	cfg     appshared.Config
	repo    ports.Repository
	storage ports.Storage
	runtime appshared.Runtime
}

func New(cfg appshared.Config, repo ports.Repository, storage ports.Storage, runtime appshared.Runtime) *UseCase {
	return &UseCase{cfg: cfg, repo: repo, storage: storage, runtime: runtime.Normalize()}
}

func (u *UseCase) Execute(ctx context.Context, cmd Command) (Result, error) {
	if cmd.UploadedBy == uuid.Nil {
		return Result{}, fileerrors.ErrUnauthenticated
	}
	if err := validate(cmd); err != nil {
		return Result{}, err
	}

	now := u.runtime.Now().UTC()
	id := u.runtime.NewUUID()
	originalName := strings.TrimSpace(cmd.OriginalName)
	mimeType := strings.TrimSpace(cmd.MimeType)
	serviceOrigin := sanitizeSegment(cmd.ServiceOrigin)
	storagePath := buildStoragePath(serviceOrigin, id, originalName, now)

	file, err := u.repo.Create(ctx, ports.FileCreate{
		ID:              id,
		OriginalName:    originalName,
		StoragePath:     storagePath,
		MimeType:        mimeType,
		SizeBytes:       cmd.SizeBytes,
		Checksum:        "",
		ServiceOrigin:   strings.TrimSpace(cmd.ServiceOrigin),
		IsPublic:        cmd.IsPublic,
		Status:          sharedv1.FileStatus_FILE_STATUS_UPLOADING,
		StorageProvider: u.cfg.StorageProvider,
		UploadedBy:      cmd.UploadedBy,
		Now:             now,
	})
	if err != nil {
		return Result{}, err
	}

	uploadURL, expiresAt, err := u.storage.PresignUpload(ctx, file.StoragePath, file.MimeType, u.cfg.UploadURLTTL)
	if err != nil {
		return Result{}, err
	}

	return Result{
		ID:        file.ID,
		UploadURL: uploadURL,
		ExpiresAt: expiresAt,
	}, nil
}

func validate(cmd Command) error {
	if strings.TrimSpace(cmd.OriginalName) == "" {
		return fileerrors.ErrInvalidInput
	}
	if strings.TrimSpace(cmd.MimeType) == "" {
		return fileerrors.ErrInvalidInput
	}
	if cmd.SizeBytes <= 0 {
		return fileerrors.ErrInvalidInput
	}
	if strings.TrimSpace(cmd.ServiceOrigin) == "" {
		return fileerrors.ErrInvalidInput
	}
	return nil
}

func buildStoragePath(serviceOrigin string, id uuid.UUID, originalName string, now time.Time) string {
	return strings.Join([]string{
		serviceOrigin,
		now.Format("2006"),
		now.Format("01"),
		now.Format("02"),
		id.String() + "-" + sanitizeFileName(originalName),
	}, "/")
}

func sanitizeFileName(value string) string {
	base := path.Base(strings.ReplaceAll(strings.TrimSpace(value), "\\", "/"))
	if base == "." || base == "/" || base == "" {
		return "file"
	}
	return sanitizeSegment(base)
}

func sanitizeSegment(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	var builder strings.Builder
	lastDash := false
	for _, char := range value {
		allowed := (char >= 'a' && char <= 'z') ||
			(char >= '0' && char <= '9') ||
			char == '.' ||
			char == '_' ||
			char == '-'
		if allowed {
			builder.WriteRune(char)
			lastDash = false
			continue
		}
		if !lastDash {
			builder.WriteByte('-')
			lastDash = true
		}
	}
	result := strings.Trim(builder.String(), "-.")
	if result == "" {
		return "file"
	}
	return result
}
