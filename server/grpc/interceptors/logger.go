package interceptors

import (
	"context"
	"github.com/charmbracelet/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"time"
)

type LoggerInterceptor struct {
	logger *log.Logger
}

func NewLoggerInterceptor(logger *log.Logger) *LoggerInterceptor {
	return &LoggerInterceptor{logger: logger}
}

func (i *LoggerInterceptor) Unary(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()

	// Извлекаем дополнительную информацию из контекста.
	clientIP := getClientIP(ctx)
	userAgent := getUserAgent(ctx)
	requestID := getRequestID(ctx)

	// Логируем начало запроса.
	i.logger.Debug("gRPC request started",
		"method", info.FullMethod,
		"client_ip", clientIP,
		"user_agent", userAgent,
		"request_id", requestID,
		"start_time", start,
	)

	// Вызываем обработчик.
	resp, err := handler(ctx, req)

	// Подготовка полей для логирования.
	duration := time.Since(start)
	statusCode := status.Code(err)
	fields := []any{
		"method", info.FullMethod,
		"client_ip", clientIP,
		"user_agent", userAgent,
		"request_id", requestID,
		"duration", duration.String(),
		"duration_ms", duration,
		"status", statusCode.String(),
		"status_code", int32(statusCode),
		"start_time", start,
		"end_time", time.Now(),
	}

	if err != nil {
		fields = append(fields,
			"err", err,
			"error_details", err.Error(),
		)
	}

	// Логируем результат запроса.
	i.logger.Log(getLogLevel(statusCode), "gRPC request completed", fields...)

	return resp, err
}

func (i *LoggerInterceptor) Stream(
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

	i.logger.Debug("gRPC stream started",
		"method", info.FullMethod,
		"client_ip", clientIP,
		"user_agent", userAgent,
		"request_id", requestID,
		"is_client_stream", info.IsClientStream,
		"is_server_stream", info.IsServerStream,
		"start_time", start,
	)

	err := handler(srv, stream)
	duration := time.Since(start)
	statusCode := status.Code(err)

	fields := []any{
		"method", info.FullMethod,
		"client_ip", clientIP,
		"user_agent", userAgent,
		"request_id", requestID,
		"is_client_stream", info.IsClientStream,
		"is_server_stream", info.IsServerStream,
		"duration", duration.String(),
		"duration_ms", duration,
		"status", statusCode.String(),
		"status_code", int32(statusCode),
		"start_time", start,
		"end_time", time.Now(),
	}

	if err != nil {
		fields = append(fields, "err", err)
	}

	i.logger.Log(getLogLevel(statusCode), "gRPC stream completed", fields...)

	return err
}

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

func getLogLevel(code codes.Code) log.Level {
	switch code {
	case codes.OK:
		return log.InfoLevel
	case codes.Canceled, codes.DeadlineExceeded:
		return log.WarnLevel
	case codes.InvalidArgument, codes.NotFound, codes.AlreadyExists,
		codes.PermissionDenied, codes.Unauthenticated:
		return log.InfoLevel
	default:
		return log.ErrorLevel
	}
}
