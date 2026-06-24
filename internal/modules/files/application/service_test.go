package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	sharedv1 "github.com/zchelalo/neuraclinic-file-management/gen/go/shared/v1"
	fileerrors "github.com/zchelalo/neuraclinic-file-management/internal/modules/files/application/errors"
	"github.com/zchelalo/neuraclinic-file-management/internal/modules/files/ports"
)

func TestRequestUploadCreatesMetadataAndPresignedURL(t *testing.T) {
	ctx := context.Background()
	fileID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	userID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	now := time.Date(2026, 6, 23, 10, 30, 0, 0, time.UTC)
	repo := newMemoryRepo()
	storage := &fakeStorage{uploadURL: "http://storage/upload", uploadExpiresAt: now.Add(15 * time.Minute)}
	service := newTestService(repo, storage, Runtime{
		Now:     func() time.Time { return now },
		NewUUID: func() uuid.UUID { return fileID },
	})

	result, err := service.RequestUpload(ctx, RequestUploadCommand{
		UploadedBy:    userID,
		OriginalName:  "clinical report.pdf",
		MimeType:      "application/pdf",
		SizeBytes:     2048,
		IsPublic:      false,
		ServiceOrigin: "record",
	})
	if err != nil {
		t.Fatalf("RequestUpload returned error: %v", err)
	}

	if result.ID != fileID {
		t.Fatalf("expected file id %s, got %s", fileID, result.ID)
	}
	if result.UploadURL != "http://storage/upload" {
		t.Fatalf("expected upload URL, got %q", result.UploadURL)
	}

	file := repo.files[fileID]
	if file.Status != sharedv1.FileStatus_FILE_STATUS_UPLOADING {
		t.Fatalf("expected uploading status, got %s", file.Status)
	}
	if file.StoragePath != "record/2026/06/23/11111111-1111-1111-1111-111111111111-clinical-report.pdf" {
		t.Fatalf("unexpected storage path: %s", file.StoragePath)
	}
	if storage.uploadKey != file.StoragePath {
		t.Fatalf("expected storage upload key %q, got %q", file.StoragePath, storage.uploadKey)
	}
}

func TestConfirmUploadRequiresObjectToExist(t *testing.T) {
	ctx := context.Background()
	fileID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	repo := newMemoryRepo()
	repo.files[fileID] = ports.File{
		ID:          fileID,
		StoragePath: "record/file.pdf",
		MimeType:    "application/pdf",
		Status:      sharedv1.FileStatus_FILE_STATUS_UPLOADING,
	}
	storage := &fakeStorage{exists: false}
	service := newTestService(repo, storage, Runtime{})

	_, err := service.ConfirmUpload(ctx, ConfirmUploadCommand{
		ID:     fileID,
		Status: sharedv1.FileStatus_FILE_STATUS_AVAILABLE,
	})
	if !errors.Is(err, ErrFailedPrecondition) {
		t.Fatalf("expected failed precondition, got %v", err)
	}
}

func TestConfirmUploadReturnsDownloadURLWhenAvailable(t *testing.T) {
	ctx := context.Background()
	fileID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	now := time.Date(2026, 6, 23, 10, 30, 0, 0, time.UTC)
	repo := newMemoryRepo()
	repo.files[fileID] = ports.File{
		ID:          fileID,
		StoragePath: "record/file.pdf",
		MimeType:    "application/pdf",
		Status:      sharedv1.FileStatus_FILE_STATUS_UPLOADING,
	}
	storage := &fakeStorage{
		exists:            true,
		downloadURL:       "http://storage/download",
		downloadExpiresAt: now.Add(10 * time.Minute),
	}
	service := newTestService(repo, storage, Runtime{Now: func() time.Time { return now }})

	result, err := service.ConfirmUpload(ctx, ConfirmUploadCommand{
		ID:     fileID,
		Status: sharedv1.FileStatus_FILE_STATUS_AVAILABLE,
	})
	if err != nil {
		t.Fatalf("ConfirmUpload returned error: %v", err)
	}

	if result.Status != sharedv1.FileStatus_FILE_STATUS_AVAILABLE {
		t.Fatalf("expected available status, got %s", result.Status)
	}
	if result.DownloadURL != "http://storage/download" {
		t.Fatalf("expected download URL, got %q", result.DownloadURL)
	}
	if repo.files[fileID].UploadedAt == nil {
		t.Fatal("expected uploaded_at to be set")
	}
}

func TestGenerateDownloadURLRequiresAvailableFile(t *testing.T) {
	ctx := context.Background()
	fileID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	repo := newMemoryRepo()
	repo.files[fileID] = ports.File{
		ID:          fileID,
		StoragePath: "record/file.pdf",
		Status:      sharedv1.FileStatus_FILE_STATUS_UPLOADING,
	}
	service := newTestService(repo, &fakeStorage{}, Runtime{})

	_, err := service.GenerateDownloadURL(ctx, fileID)
	if !errors.Is(err, ErrFailedPrecondition) {
		t.Fatalf("expected failed precondition, got %v", err)
	}
}

func newTestService(repo ports.Repository, storage ports.Storage, runtime Runtime) *Service {
	return NewServiceWithRuntime(Config{
		StorageProvider: "s3",
		UploadURLTTL:    15 * time.Minute,
		DownloadURLTTL:  10 * time.Minute,
	}, repo, storage, runtime)
}

type memoryRepo struct {
	files map[uuid.UUID]ports.File
}

func newMemoryRepo() *memoryRepo {
	return &memoryRepo{files: make(map[uuid.UUID]ports.File)}
}

func (r *memoryRepo) Create(_ context.Context, file ports.FileCreate) (ports.File, error) {
	created := ports.File{
		ID:              file.ID,
		OriginalName:    file.OriginalName,
		StoragePath:     file.StoragePath,
		MimeType:        file.MimeType,
		SizeBytes:       file.SizeBytes,
		Checksum:        file.Checksum,
		ServiceOrigin:   file.ServiceOrigin,
		IsPublic:        file.IsPublic,
		Status:          file.Status,
		StorageProvider: file.StorageProvider,
		UploadedBy:      file.UploadedBy,
		CreatedAt:       file.Now,
		UpdatedAt:       file.Now,
	}
	r.files[created.ID] = created
	return created, nil
}

func (r *memoryRepo) ByID(_ context.Context, id uuid.UUID) (ports.File, error) {
	file, ok := r.files[id]
	if !ok {
		return ports.File{}, fileerrors.ErrNotFound
	}
	return file, nil
}

func (r *memoryRepo) ConfirmUpload(_ context.Context, id uuid.UUID, status sharedv1.FileStatus, now time.Time) (ports.File, error) {
	file, ok := r.files[id]
	if !ok {
		return ports.File{}, fileerrors.ErrNotFound
	}
	file.Status = status
	file.UpdatedAt = now
	if status == sharedv1.FileStatus_FILE_STATUS_AVAILABLE {
		file.UploadedAt = &now
	}
	r.files[id] = file
	return file, nil
}

type fakeStorage struct {
	uploadURL         string
	uploadExpiresAt   time.Time
	uploadKey         string
	downloadURL       string
	downloadExpiresAt time.Time
	exists            bool
	existsErr         error
}

func (s *fakeStorage) PresignUpload(_ context.Context, key, _ string, _ time.Duration) (string, time.Time, error) {
	s.uploadKey = key
	return s.uploadURL, s.uploadExpiresAt, nil
}

func (s *fakeStorage) PresignDownload(_ context.Context, _ string, _ time.Duration) (string, time.Time, error) {
	return s.downloadURL, s.downloadExpiresAt, nil
}

func (s *fakeStorage) Exists(_ context.Context, _ string) (bool, error) {
	return s.exists, s.existsErr
}

var _ ports.Repository = (*memoryRepo)(nil)
var _ ports.Storage = (*fakeStorage)(nil)
