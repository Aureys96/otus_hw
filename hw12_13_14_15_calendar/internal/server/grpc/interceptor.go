package grpc

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func logInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	h, err := handler(ctx, req)

	zap.L().Info("grpc request",
		zap.String("method", info.FullMethod),
		zap.Duration("duration", time.Since(start)),
		zap.Error(err),
	)

	return h, err
}

func logRequest(logger *zap.Logger) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		h, err := handler(ctx, req)

		logger.Info("grpc request",
			zap.String("method", info.FullMethod),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)

		return h, err
	}
}
