package ctld

import (
	"context"

	commonv1 "master-otel/internal/proto/common/v1"
	"master-otel/pkg/log"

	"go.uber.org/zap"
)

func (s *Service) CreateUser(ctx context.Context, req *commonv1.User) (*commonv1.User, error) {
	log.WithCtx(ctx).Info("create user", zap.String("username", req.GetUsername()))
	return s.storedClient.CreateUser(ctx, req)
}
