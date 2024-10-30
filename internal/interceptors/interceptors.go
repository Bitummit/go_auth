package interceptors

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
)

func UnaryLogRequest(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		log.Info("request: %v", slog.Attr{
			Key: "info",
			Value: slog.StringValue(info.FullMethod),
		})

		return handler(ctx, req)
	}
}