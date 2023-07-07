package utils

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"master-otel/pkg/log"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	KeyTraceID     = "trace_id"
	KeyServiceName = "service_name"
)

// GrpcDial dials to the gRPC server
func GrpcDial(ctx context.Context, addr string, dialOpts ...grpc.DialOption) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock())
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, dialOpts...)
	conn, err := grpc.DialContext(ctx, addr, opts...)
	if err != nil {
		return nil, fmt.Errorf("dial %s: %w", addr, err)
	}
	return conn, nil
}

func TraceUnaryServerInterceptor(serivce string) grpc.UnaryServerInterceptor {
	return grpc.UnaryServerInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		incoming, ok := metadata.FromIncomingContext(ctx)
		if ok {
			trace := uuid.Nil.String()
			if value := incoming.Get(KeyTraceID); len(value) > 0 {
				trace = value[0]
			}
			ctx = log.AddCtx(ctx, zap.String(KeyServiceName, serivce), zap.Any(KeyTraceID, trace))
			ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(KeyTraceID, trace))
		}

		return handler(ctx, req)
	})
}

func LoggerUnaryServerInterceptor() logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, level logging.Level, msg string, pairs ...any) {
		if os.Getenv("GRPC_LOG_DEBUG") == "0" {
			return
		}
		if msg == "started call" {
			return
		}

		var fields []zapcore.Field
		for i := 0; i < len(pairs); i += 2 {
			key, value := pairs[i].(string), pairs[i+1].(string)
			switch key {
			case "grpc.method":
				fields = append(fields, zap.String("method", value))
			case "peer.address":
				fields = append(fields, zap.String("peer", value))
			case "grpc.time_ms":
				fields = append(fields, zap.String("elapsed", value+"ms"))
			case "grpc.code":
				fields = append(fields, zap.String("code", value))
			case "grpc.error":
				fields = append(fields, zap.String("error", value))
			}
		}
		content := "grpc"
		switch level {
		case logging.LevelDebug:
			log.WithCtx(ctx).Debug(content, fields...)
		case logging.LevelInfo:
			log.WithCtx(ctx).Info(content, fields...)
		case logging.LevelWarn:
			log.WithCtx(ctx).Warn(content, fields...)
		case logging.LevelError:
			log.WithCtx(ctx).Error(content, fields...)
		}
	})
}

func TraceMiddleware(service string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			trace := NewUUID()
			ctx := log.AddCtx(c.Request().Context(), zap.String(KeyServiceName, service), zap.String(KeyTraceID, trace))
			ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(KeyTraceID, trace))
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

func NewUUID() string {
	return uuid.NewString()[0:8]
}

func ContextWithTrace(ctx context.Context, service string) context.Context {
	trace := NewUUID()
	return log.AddCtx(ctx, zap.String(KeyServiceName, service), zap.String(KeyTraceID, trace))
}

func LoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			realIp := c.RealIP()
			method := c.Request().Method
			uri := c.Request().RequestURI
			log.WithCtx(c.Request().Context()).Info(
				"http",
				zap.String("real_ip", realIp),
				zap.String("uri", uri),
				zap.String("method", method),
				zap.Int64("bytes_in", c.Request().ContentLength),
			)
			err := next(c)
			status := c.Response().Status
			elapsed := time.Since(start)
			size := c.Response().Size
			log.WithCtx(c.Request().Context()).Info(
				"http",
				zap.String("real_ip", realIp),
				zap.String("uri", uri),
				zap.String("method", method),
				zap.Int64("bytes_out", size),
				zap.Int("code", status),
				zap.Error(err),
				zap.String("elapsed", elapsed.String()),
			)
			return nil
		}
	}
}

func NewGrpcServer(service string) *grpc.Server {
	return grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(func(p any) (err error) {
				log.Error("recovered from panic", zap.Any("panic", p), zap.Any("stack", debug.Stack()))
				return status.Errorf(codes.Internal, "%s", p)
			})),
			grpc.UnaryServerInterceptor(TraceUnaryServerInterceptor(service)),
			logging.UnaryServerInterceptor(LoggerUnaryServerInterceptor()),
		),
		grpc.ChainStreamInterceptor(
			recovery.StreamServerInterceptor(recovery.WithRecoveryHandler(func(p any) (err error) {
				log.Error("recovered from panic", zap.Any("panic", p), zap.Any("stack", debug.Stack()))
				return status.Errorf(codes.Internal, "%s", p)
			})),
		),
	)
}
