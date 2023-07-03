package ctld

import (
	"context"

	commonv1 "master-otel/internal/proto/common/v1"
	"master-otel/pkg/log"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/uptrace/opentelemetry-go-extra/otelplay"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func (s *Service) CreateUser(ctx context.Context, req *commonv1.User) (*commonv1.User, error) {
	log.WithContext(ctx).Info("create user", zap.String("trace", otelplay.TraceURL(trace.SpanFromContext(ctx))))
	return s.storedClient.CreateUser(ctx, req)
}

func (s *Service) GetUser(ctx context.Context, req *commonv1.Identity) (*commonv1.User, error) {
	return s.storedClient.GetUser(ctx, req)
}

func (s *Service) DeleteUser(ctx context.Context, req *commonv1.Identity) (*empty.Empty, error) {
	return s.storedClient.DeleteUser(ctx, req)
}
