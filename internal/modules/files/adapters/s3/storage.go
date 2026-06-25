package s3storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/zchelalo/neuraclinic-file-management/internal/modules/files/ports"
)

type Config struct {
	Bucket          string
	Region          string
	Endpoint        string
	PublicEndpoint  string
	AccessKeyID     string
	SecretAccessKey string
	ForcePathStyle  bool
}

type Storage struct {
	bucket    string
	client    *s3.Client
	presigner *s3.PresignClient
}

func New(ctx context.Context, cfg Config) (*Storage, error) {
	if cfg.Bucket == "" {
		return nil, fmt.Errorf("bucket is required")
	}
	if cfg.Region == "" {
		return nil, fmt.Errorf("region is required")
	}

	loadOptions := []func(*config.LoadOptions) error{
		config.WithRegion(cfg.Region),
	}
	if cfg.AccessKeyID != "" || cfg.SecretAccessKey != "" {
		loadOptions = append(loadOptions, config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		))
	}

	awsCfg, err := config.LoadDefaultConfig(ctx, loadOptions...)
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg, func(options *s3.Options) {
		options.UsePathStyle = cfg.ForcePathStyle
		if cfg.Endpoint != "" {
			options.BaseEndpoint = aws.String(cfg.Endpoint)
		}
	})

	presignClient := client
	if cfg.PublicEndpoint != "" {
		presignClient = s3.NewFromConfig(awsCfg, func(options *s3.Options) {
			options.UsePathStyle = cfg.ForcePathStyle
			options.BaseEndpoint = aws.String(cfg.PublicEndpoint)
		})
	}

	return &Storage{
		bucket:    cfg.Bucket,
		client:    client,
		presigner: s3.NewPresignClient(presignClient),
	}, nil
}

func (s *Storage) PresignUpload(ctx context.Context, key, contentType string, expires time.Duration) (string, time.Time, error) {
	result, err := s.presigner.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, func(options *s3.PresignOptions) {
		options.Expires = expires
	})
	if err != nil {
		return "", time.Time{}, fmt.Errorf("presign put object: %w", err)
	}

	return result.URL, time.Now().UTC().Add(expires), nil
}

func (s *Storage) PresignDownload(ctx context.Context, key string, expires time.Duration) (string, time.Time, error) {
	result, err := s.presigner.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}, func(options *s3.PresignOptions) {
		options.Expires = expires
	})
	if err != nil {
		return "", time.Time{}, fmt.Errorf("presign get object: %w", err)
	}

	return result.URL, time.Now().UTC().Add(expires), nil
}

func (s *Storage) Exists(ctx context.Context, key string) (bool, error) {
	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err == nil {
		return true, nil
	}

	var notFound *types.NotFound
	if errors.As(err, &notFound) {
		return false, nil
	}

	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.ErrorCode() {
		case "NotFound", "NoSuchKey", "404":
			return false, nil
		}
	}

	return false, fmt.Errorf("head object: %w", err)
}

var _ ports.Storage = (*Storage)(nil)
