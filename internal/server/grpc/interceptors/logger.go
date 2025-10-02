package interceptors

import (
	"context"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func UnaryLogger(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		// Извлекаем дополнительную информацию из контекста
		clientIP := getClientIP(ctx)
		userAgent := getUserAgent(ctx)
		requestID := getRequestID(ctx)

		// Логируем начало запроса
		logger.Debug("gRPC request started",
			zap.String("method", info.FullMethod),
			zap.String("client_ip", clientIP),
			zap.String("user_agent", userAgent),
			zap.String("request_id", requestID),
			zap.Time("start_time", start),
		)

		// Вызываем обработчик
		resp, err := handler(ctx, req)

		// Подготовка полей для логирования
		duration := time.Since(start)
		statusCode := status.Code(err)

		fields := []zap.Field{
			zap.String("method", info.FullMethod),
			zap.String("client_ip", clientIP),
			zap.String("user_agent", userAgent),
			zap.String("request_id", requestID),
			zap.String("duration", duration.String()),
			zap.Duration("duration_ms", duration),
			zap.String("status", statusCode.String()),
			zap.Int32("status_code", int32(statusCode)),
			zap.Time("start_time", start),
			zap.Time("end_time", time.Now()),
		}

		// Логируем в зависимости от статуса
		logLevel := getLogLevel(statusCode)

		if err != nil {
			fields = append(fields,
				zap.Error(err),
				zap.String("error_details", err.Error()),
			)
		}

		// Логируем с соответствующим уровнем
		logger.Log(logLevel, "gRPC request completed", fields...)

		return resp, err
	}
}

// StreamInterceptor для streaming RPC
func StreamLogger(logger *zap.Logger) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		start := time.Now()

		ctx := stream.Context()
		clientIP := getClientIP(ctx)
		userAgent := getUserAgent(ctx)
		requestID := getRequestID(ctx)

		logger.Debug("gRPC stream started",
			zap.String("method", info.FullMethod),
			zap.String("client_ip", clientIP),
			zap.String("user_agent", userAgent),
			zap.String("request_id", requestID),
			zap.Bool("is_client_stream", info.IsClientStream),
			zap.Bool("is_server_stream", info.IsServerStream),
			zap.Time("start_time", start),
		)

		err := handler(srv, stream)
		duration := time.Since(start)
		statusCode := status.Code(err)

		fields := []zap.Field{
			zap.String("method", info.FullMethod),
			zap.String("client_ip", clientIP),
			zap.String("user_agent", userAgent),
			zap.String("request_id", requestID),
			zap.Bool("is_client_stream", info.IsClientStream),
			zap.Bool("is_server_stream", info.IsServerStream),
			zap.String("duration", duration.String()),
			zap.Duration("duration_ms", duration),
			zap.String("status", statusCode.String()),
			zap.Int32("status_code", int32(statusCode)),
			zap.Time("start_time", start),
			zap.Time("end_time", time.Now()),
		}

		logLevel := getLogLevel(statusCode)
		if err != nil {
			fields = append(fields, zap.Error(err))
		}

		logger.Log(logLevel, "gRPC stream completed", fields...)

		return err
	}
}

// Вспомогательные функции
func getClientIP(ctx context.Context) string {
	if p, ok := peer.FromContext(ctx); ok {
		return p.Addr.String()
	}
	return "unknown"
}

func getUserAgent(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if userAgents := md.Get("user-agent"); len(userAgents) > 0 {
			return userAgents[0]
		}
	}
	return "unknown"
}

func getRequestID(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if requestIDs := md.Get("x-request-id"); len(requestIDs) > 0 {
			return requestIDs[0]
		}
		if requestIDs := md.Get("request-id"); len(requestIDs) > 0 {
			return requestIDs[0]
		}
	}
	return "unknown"
}

func getLogLevel(code codes.Code) zapcore.Level {
	switch code {
	case codes.OK:
		return zapcore.InfoLevel
	case codes.Canceled, codes.DeadlineExceeded:
		return zapcore.WarnLevel
	case codes.InvalidArgument, codes.NotFound, codes.AlreadyExists,
		codes.PermissionDenied, codes.Unauthenticated:
		return zapcore.InfoLevel
	default:
		return zapcore.ErrorLevel
	}
}
