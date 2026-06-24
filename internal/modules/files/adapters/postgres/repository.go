package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	sharedv1 "github.com/zchelalo/neuraclinic-file-management/gen/go/shared/v1"
	filesdb "github.com/zchelalo/neuraclinic-file-management/internal/db/sqlc/files"
	fileerrors "github.com/zchelalo/neuraclinic-file-management/internal/modules/files/application/errors"
	"github.com/zchelalo/neuraclinic-file-management/internal/modules/files/ports"
	pgutil "github.com/zchelalo/neuraclinic-file-management/internal/shared/postgresutil"
)

type Repository struct {
	q *filesdb.Queries
}

func NewRepository(db filesdb.DBTX) *Repository {
	return &Repository{q: filesdb.New(db)}
}

func (r *Repository) Create(ctx context.Context, file ports.FileCreate) (ports.File, error) {
	row, err := r.q.CreateFile(ctx, filesdb.CreateFileParams{
		ID:              pgutil.UUID(file.ID),
		OriginalName:    file.OriginalName,
		StoragePath:     file.StoragePath,
		MimeType:        file.MimeType,
		SizeBytes:       file.SizeBytes,
		Checksum:        file.Checksum,
		ServiceOrigin:   file.ServiceOrigin,
		IsPublic:        file.IsPublic,
		Status:          file.Status.String(),
		StorageProvider: file.StorageProvider,
		UploadedBy:      pgutil.UUID(file.UploadedBy),
		CreatedAt:       pgutil.Timestamptz(file.Now.UTC()),
	})
	if err != nil {
		return ports.File{}, err
	}
	return fileFromRow(row), nil
}

func (r *Repository) ByID(ctx context.Context, id uuid.UUID) (ports.File, error) {
	row, err := r.q.GetFileByID(ctx, pgutil.UUID(id))
	if err != nil {
		return ports.File{}, mapNoRows(err)
	}
	return fileFromRow(row), nil
}

func (r *Repository) ConfirmUpload(ctx context.Context, id uuid.UUID, status sharedv1.FileStatus, now time.Time) (ports.File, error) {
	row, err := r.q.ConfirmFileUpload(ctx, filesdb.ConfirmFileUploadParams{
		ID:        pgutil.UUID(id),
		Status:    status.String(),
		UpdatedAt: pgutil.Timestamptz(now.UTC()),
	})
	if err != nil {
		return ports.File{}, mapNoRows(err)
	}
	return fileFromRow(row), nil
}

func fileFromRow(row filesdb.File) ports.File {
	return ports.File{
		ID:              pgutil.UUIDValue(row.ID),
		OriginalName:    row.OriginalName,
		StoragePath:     row.StoragePath,
		MimeType:        row.MimeType,
		SizeBytes:       row.SizeBytes,
		Checksum:        row.Checksum,
		ServiceOrigin:   row.ServiceOrigin,
		IsPublic:        row.IsPublic,
		Status:          parseFileStatus(row.Status),
		StorageProvider: row.StorageProvider,
		UploadedBy:      pgutil.UUIDValue(row.UploadedBy),
		UploadedAt:      pgutil.TimestamptzPtr(row.UploadedAt),
		CreatedAt:       row.CreatedAt.Time,
		UpdatedAt:       row.UpdatedAt.Time,
		DeletedAt:       pgutil.TimestamptzPtr(row.DeletedAt),
	}
}

func parseFileStatus(value string) sharedv1.FileStatus {
	if parsed, ok := sharedv1.FileStatus_value[value]; ok {
		return sharedv1.FileStatus(parsed)
	}
	return sharedv1.FileStatus_FILE_STATUS_UNSPECIFIED
}

func mapNoRows(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return fileerrors.ErrNotFound
	}
	return err
}

var _ ports.Repository = (*Repository)(nil)
