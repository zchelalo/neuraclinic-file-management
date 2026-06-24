package bootstrap

import (
	"context"
	"fmt"

	filesgrpc "github.com/zchelalo/neuraclinic-file-management/internal/modules/files/adapters/grpc"
	filespg "github.com/zchelalo/neuraclinic-file-management/internal/modules/files/adapters/postgres"
	s3storage "github.com/zchelalo/neuraclinic-file-management/internal/modules/files/adapters/s3"
	"github.com/zchelalo/neuraclinic-file-management/internal/modules/files/application"
	grpcserver "github.com/zchelalo/neuraclinic-file-management/internal/server/grpc"
	"go.uber.org/zap"
)

type App struct {
	Server  *grpcserver.Server
	Cleanup func(context.Context) error
}

func InitApp(ctx context.Context, logger *zap.Logger, cfg Config) (*App, error) {
	db, err := NewDB(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("cannot initialize db: %w", err)
	}

	storage, err := s3storage.New(ctx, s3storage.Config{
		Bucket:          cfg.StorageBucket,
		Region:          cfg.StorageRegion,
		Endpoint:        cfg.StorageEndpoint,
		PublicEndpoint:  cfg.StoragePublicEndpoint,
		AccessKeyID:     cfg.StorageAccessKeyID,
		SecretAccessKey: cfg.StorageSecretAccessKey,
		ForcePathStyle:  cfg.StorageForcePathStyle,
	})
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("cannot initialize storage: %w", err)
	}

	repo := filespg.NewRepository(db)
	filesApp := application.NewService(application.Config{
		StorageProvider: cfg.StorageProvider,
		UploadURLTTL:    cfg.StorageUploadURLTTL,
		DownloadURLTTL:  cfg.StorageDownloadURLTTL,
	}, repo, storage)

	server, err := grpcserver.New(grpcserver.Config{
		Port:            cfg.Port,
		ServiceName:     cfg.ServiceName,
		TLSCertFilePath: cfg.GRPCTLSCertPath,
		TLSKeyFilePath:  cfg.GRPCTLSKeyPath,
	}, logger, filesgrpc.NewServices(filesApp))
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("cannot create grpc server: %w", err)
	}

	return &App{
		Server: server,
		Cleanup: func(context.Context) error {
			server.GracefulStop()
			db.Close()
			return nil
		},
	}, nil
}
