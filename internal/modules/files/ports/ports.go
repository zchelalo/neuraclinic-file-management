package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
	sharedv1 "github.com/zchelalo/neuraclinic-file-management/gen/go/shared/v1"
)

type File struct {
	ID              uuid.UUID
	OriginalName    string
	StoragePath     string
	MimeType        string
	SizeBytes       int64
	Checksum        string
	ServiceOrigin   string
	IsPublic        bool
	Status          sharedv1.FileStatus
	StorageProvider string
	UploadedBy      uuid.UUID
	UploadedAt      *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}

type FileCreate struct {
	ID              uuid.UUID
	OriginalName    string
	StoragePath     string
	MimeType        string
	SizeBytes       int64
	Checksum        string
	ServiceOrigin   string
	IsPublic        bool
	Status          sharedv1.FileStatus
	StorageProvider string
	UploadedBy      uuid.UUID
	Now             time.Time
}

type FileStatusChangedEvent struct {
	EventID       uuid.UUID
	FileID        uuid.UUID
	ServiceOrigin string
	Status        sharedv1.FileStatus
	OccurredAt    time.Time
	RequestID     string
	TraceID       string
}

type Repository interface {
	Create(ctx context.Context, file FileCreate) (File, error)
	ByID(ctx context.Context, id uuid.UUID) (File, error)
	ConfirmUpload(ctx context.Context, id uuid.UUID, status sharedv1.FileStatus, now time.Time) (File, error)
}

type Storage interface {
	PresignUpload(ctx context.Context, key, contentType string, expires time.Duration) (string, time.Time, error)
	PresignDownload(ctx context.Context, key string, expires time.Duration) (string, time.Time, error)
	Exists(ctx context.Context, key string) (bool, error)
}

type EventPublisher interface {
	PublishFileStatusChanged(ctx context.Context, event FileStatusChangedEvent) error
	Close() error
}
