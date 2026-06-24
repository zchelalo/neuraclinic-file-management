package grpc

import (
	filemanagementv1 "github.com/zchelalo/neuraclinic-file-management/gen/go/file_management/v1"
	"github.com/zchelalo/neuraclinic-file-management/internal/modules/files/application"
)

type FileManagementService struct {
	filemanagementv1.UnimplementedFileManagementServiceServer
	app *application.Service
}

type Services struct {
	FileManagement *FileManagementService
}

func NewServices(app *application.Service) *Services {
	return &Services{FileManagement: NewFileManagementService(app)}
}

func NewFileManagementService(app *application.Service) *FileManagementService {
	return &FileManagementService{app: app}
}
