package grpc

import (
	"context"

	filemanagementv1 "github.com/zchelalo/neuraclinic-file-management/gen/go/file_management/v1"
	"github.com/zchelalo/neuraclinic-file-management/internal/modules/files/application"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *FileManagementService) RequestUpload(ctx context.Context, req *filemanagementv1.FileManagementServiceRequestUploadRequest) (*filemanagementv1.FileManagementServiceRequestUploadResponse, error) {
	uploadedBy, err := userID(ctx)
	if err != nil {
		return nil, mapError(err)
	}

	result, err := s.app.RequestUpload(ctx, application.RequestUploadCommand{
		UploadedBy:    uploadedBy,
		OriginalName:  req.GetOriginalName(),
		MimeType:      req.GetMimeType(),
		SizeBytes:     req.GetSizeBytes(),
		IsPublic:      req.GetIsPublic(),
		ServiceOrigin: req.GetServiceOrigin(),
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &filemanagementv1.FileManagementServiceRequestUploadResponse{
		Id:        result.ID.String(),
		UploadUrl: result.UploadURL,
		ExpiresAt: timestamppb.New(result.ExpiresAt),
	}, nil
}

func (s *FileManagementService) ConfirmUpload(ctx context.Context, req *filemanagementv1.FileManagementServiceConfirmUploadRequest) (*filemanagementv1.FileManagementServiceConfirmUploadResponse, error) {
	id, err := parseID(req.GetId())
	if err != nil {
		return nil, mapError(err)
	}

	result, err := s.app.ConfirmUpload(ctx, application.ConfirmUploadCommand{
		ID:     id,
		Status: req.GetStatus(),
	})
	if err != nil {
		return nil, mapError(err)
	}

	response := &filemanagementv1.FileManagementServiceConfirmUploadResponse{
		Id:     result.ID.String(),
		Status: result.Status,
	}
	if result.DownloadURL != "" {
		response.DownloadUrl = &result.DownloadURL
		response.ExpiresAt = timestamppb.New(result.ExpiresAt)
	}

	return response, nil
}

func (s *FileManagementService) GenerateDownloadUrl(ctx context.Context, req *filemanagementv1.FileManagementServiceGenerateDownloadUrlRequest) (*filemanagementv1.FileManagementServiceGenerateDownloadUrlResponse, error) {
	id, err := parseID(req.GetId())
	if err != nil {
		return nil, mapError(err)
	}

	result, err := s.app.GenerateDownloadURL(ctx, id)
	if err != nil {
		return nil, mapError(err)
	}

	return &filemanagementv1.FileManagementServiceGenerateDownloadUrlResponse{
		DownloadUrl: result.DownloadURL,
		ExpiresAt:   timestamppb.New(result.ExpiresAt),
	}, nil
}
