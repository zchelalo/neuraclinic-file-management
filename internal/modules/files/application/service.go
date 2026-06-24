package application

import (
	"context"

	"github.com/google/uuid"
	sharedv1 "github.com/zchelalo/neuraclinic-file-management/gen/go/shared/v1"
	"github.com/zchelalo/neuraclinic-file-management/internal/modules/files/application/confirmupload"
	"github.com/zchelalo/neuraclinic-file-management/internal/modules/files/application/generatedownloadurl"
	"github.com/zchelalo/neuraclinic-file-management/internal/modules/files/application/requestupload"
	appshared "github.com/zchelalo/neuraclinic-file-management/internal/modules/files/application/shared"
	"github.com/zchelalo/neuraclinic-file-management/internal/modules/files/ports"
)

type Config = appshared.Config
type Runtime = appshared.Runtime

type Service struct {
	requestUpload       *requestupload.UseCase
	confirmUpload       *confirmupload.UseCase
	generateDownloadURL *generatedownloadurl.UseCase
}

func NewService(cfg Config, repo ports.Repository, storage ports.Storage) *Service {
	return NewServiceWithRuntime(cfg, repo, storage, appshared.DefaultRuntime())
}

func NewServiceWithRuntime(cfg Config, repo ports.Repository, storage ports.Storage, runtime Runtime) *Service {
	runtime = runtime.Normalize()
	return &Service{
		requestUpload:       requestupload.New(cfg, repo, storage, runtime),
		confirmUpload:       confirmupload.New(cfg, repo, storage, runtime),
		generateDownloadURL: generatedownloadurl.New(cfg, repo, storage),
	}
}

func DefaultRuntime() Runtime {
	return appshared.DefaultRuntime()
}

func (s *Service) RequestUpload(ctx context.Context, cmd requestupload.Command) (requestupload.Result, error) {
	return s.requestUpload.Execute(ctx, cmd)
}

func (s *Service) ConfirmUpload(ctx context.Context, cmd confirmupload.Command) (confirmupload.Result, error) {
	return s.confirmUpload.Execute(ctx, cmd)
}

func (s *Service) GenerateDownloadURL(ctx context.Context, id uuid.UUID) (generatedownloadurl.Result, error) {
	return s.generateDownloadURL.Execute(ctx, generatedownloadurl.Command{ID: id})
}

type RequestUploadCommand = requestupload.Command
type RequestUploadResult = requestupload.Result
type ConfirmUploadCommand = confirmupload.Command
type ConfirmUploadResult = confirmupload.Result
type GenerateDownloadURLResult = generatedownloadurl.Result

var _ = sharedv1.FileStatus_FILE_STATUS_UNSPECIFIED
