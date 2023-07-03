package utils

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"master-otel/pkg/log"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// GrpcDial dials to the gRPC server
func GrpcDial(ctx context.Context, addr string, dialOpts ...grpc.DialOption) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock())
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()))
	opts = append(opts, dialOpts...)
	conn, err := grpc.DialContext(ctx, addr, opts...)
	if err != nil {
		return nil, fmt.Errorf("dial %s: %w", addr, err)
	}
	return conn, nil
}

func NewGrpcServer() *grpc.Server {
	return grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc.UnaryServerInterceptor(otelgrpc.UnaryServerInterceptor()),
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
