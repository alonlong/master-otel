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
	e := models.User{
		Email:    req.GetEmail(),
		Username: req.GetUsername(),
	}
	if err := s.store.CreateUser(ctx, &e); err != nil {
		return nil, status.Errorf(codes.Internal, "create user %s: %v", req.GetEmail(), err)
	}
	return e.ToProto(), nil
}

func (s *Service) DeleteUser(ctx context.Context, req *commonv1.Identity) (*commonv1.Empty, error) {
	log.WithCtx(ctx).Info("delete user", zap.Int64("id", req.GetId()))
	user, err := s.store.GetUser(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get user %d: %v", req.GetId(), err)
	}
	if user == nil {
		return nil, status.Errorf(codes.NotFound, "user %d not found", req.GetId())
	}
	if err := s.store.DeleteUser(ctx, req.GetId()); err != nil {
		return nil, status.Errorf(codes.Internal, "delete user %d: %v", req.GetId(), err)
	}
	return &commonv1.Empty{}, nil
}
