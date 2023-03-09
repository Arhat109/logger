package grpc

import (
	"context"
	"github.com/Arhat109/logger/pkg/dto"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"time"
)

// QueryAndSetTracingId -- проверяет наличие уникального ид запроса в контексте, и устанавливает при отсутствии
func QueryAndSetTracingId(ctx context.Context, name string) {
	tag := ctx.Value(name)
	if tag == nil {
		if u, err := uuid.NewRandom(); err == nil {
			ctx = context.WithValue(ctx, name, u.String())
		}
	}
}

// UnaryTracingInterceptor -- отслеживает наличие уникального идента запроса и создает его в случае отсутствия
func UnaryTracingInterceptor(lgr dto.Loggable) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		QueryAndSetTracingId(ctx, lgr.GetTraceId())
		return handler(ctx, req)
	}
}

// UnaryLoggerInterceptor returns a new unary server interceptors that adds zap.Logger to the context.
func UnaryLoggerInterceptor(lgr dto.Loggable) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		startedAt := time.Now()
		resp, err := handler(ctx, req)
		executed := time.Since(startedAt)
		if err != nil {
			lgr.ErrorCtx(ctx, "%s, executedTime=%d usec.", err.Error(), executed.Microseconds())
		} else {
			lgr.InfoCtx(ctx, "UnaryLoggerInterceptor(): success from %s, executedTime=%d usec.", info.FullMethod, executed.Microseconds())
		}
		return resp, err
	}
}

// StreamTracingInterceptor -- отслеживает наличие уникального идента запроса и создает его в случае отсутствия
func StreamTracingInterceptor(lgr dto.Loggable) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		QueryAndSetTracingId(stream.Context(), lgr.GetTraceId())
		return handler(srv, stream)
	}
}

// StreamLoggerInterceptor returns a new unary server interceptors that adds zap.Logger to the context.
func StreamLoggerInterceptor(lgr dto.Loggable) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		startedAt := time.Now()
		err := handler(srv, stream)
		executed := time.Since(startedAt)
		if err != nil {
			lgr.ErrorCtx(stream.Context(), ", executedTime=%d usec.", err.Error(), executed.Microseconds())
		} else {
			lgr.InfoCtx(stream.Context(), "StreamLoggerInterceptor(): success from %s, executedTime=%d usec.", info.FullMethod, executed.Microseconds())
		}
		return err
	}
}
