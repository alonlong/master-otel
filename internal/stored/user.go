package stored

import (
	"context"

	"master-otel/internal/entity/models"
	commonv1 "master-otel/internal/proto/common/v1"
	"master-otel/pkg/log"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) CreateUser(ctx context.Context, req *commonv1.User) (*commonv1.User, error) {
	log.WithCtx(ctx).Info("create user", zap.String("username", req.GetUsername()))
	entity := models.User{
		Email:    req.GetEmail(),
		Username: req.GetUsername(),
	}
	if err := s.store.CreateUser(ctx, &entity); err != nil {
		return nil, status.Errorf(codes.Internal, "create user %s: %v", req.GetEmail(), err)
	}
	return entity.ToProto(), nil
}
