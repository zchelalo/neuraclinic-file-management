package grpcserver

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/zchelalo/neuraclinic-file-management/internal/shared/appctx"
	"github.com/zchelalo/neuraclinic-file-management/internal/shared/uuidx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	headerRequestID = "x-request-id"
	headerTraceID   = "x-trace-id"
	headerUserID    = "x-user-id"
)

func UnaryInterceptor(baseLogger *zap.Logger, serviceName string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		startedAt := time.Now()
		requestID := metadataValue(ctx, headerRequestID)
		if requestID == "" {
			requestID = uuidx.NewString()
		}
		traceID := metadataValue(ctx, headerTraceID)
		if traceID == "" {
			traceID = uuidx.NewString()
		}

		logger := baseLogger.With(
			zap.String("service", serviceName),
			zap.String("request_id", requestID),
			zap.String("trace_id", traceID),
			zap.String("grpc_method", info.FullMethod),
		)

		ctx = appctx.WithLogger(ctx, logger)
		ctx = appctx.WithRequestID(ctx, requestID)
		ctx = appctx.WithTraceID(ctx, traceID)
		ctx = withOptionalUUID(ctx, headerUserID, appctx.WithUserID)

		defer func() {
			if recovered := recover(); recovered != nil {
				logger.Error("panic recovered", zap.Any("panic", recovered))
				err = status.Error(codes.Internal, "internal server error")
			}

			code := status.Code(err)
			fields := []zap.Field{
				zap.String("grpc_code", code.String()),
				zap.Duration("duration", time.Since(startedAt)),
			}
			if err != nil {
				fields = append(fields, zap.Error(err))
			}

			if code == codes.OK {
				logger.Info("grpc request completed", fields...)
				return
			}
			logger.Warn("grpc request failed", fields...)
		}()

		return handler(ctx, req)
	}
}

func metadataValue(ctx context.Context, key string) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	values := md.Get(key)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func withOptionalUUID(ctx context.Context, header string, setter func(context.Context, uuid.UUID) context.Context) context.Context {
	value := metadataValue(ctx, header)
	if value == "" {
		return ctx
	}
	id, err := uuid.Parse(value)
	if err != nil || id == uuid.Nil {
		return ctx
	}
	return setter(ctx, id)
}
