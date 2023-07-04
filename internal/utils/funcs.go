package utils

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"master-otel/pkg/log"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	KeyTraceID     = "trace-id"
	KeyServiceName = "service"
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

func TraceMiddleware(service string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			trace := uuid.New().String()
			ctx := log.AddCtx(c.Request().Context(), zap.String(KeyServiceName, service), zap.String(KeyTraceID, trace))
			ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(KeyTraceID, trace))
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

func NewGrpcServer(service string) *grpc.Server {
	return grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc.UnaryServerInterceptor(TraceUnaryServerInterceptor(service)),
			recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(func(p any) (err error) {
				log.Error("recovered from panic", zap.Any("panic", p), zap.Any("stack", debug.Stack()))
				return status.Errorf(codes.Internal, "%s", p)
			})),
		),
		grpc.ChainStreamInterceptor(
			recovery.StreamServerInterceptor(recovery.WithRecoveryHandler(func(p any) (err error) {
				log.Error("recovered from panic", zap.Any("panic", p), zap.Any("stack", debug.Stack()))
				return status.Errorf(codes.Internal, "%s", p)
			})),
		),
	)
}
