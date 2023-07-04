package stored

import (
	"context"

	"master-otel/internal/entity/models"
	commonv1 "master-otel/internal/proto/common/v1"
	"master-otel/pkg/log"

	"github.com/golang/protobuf/ptypes/empty"
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

func (s *Service) GetUser(ctx context.Context, req *commonv1.Identity) (*commonv1.User, error) {
	entity, err := s.store.GetUser(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "query user %d: %v", req.GetId(), err)
	}
	if entity == nil {
		return nil, status.Errorf(codes.NotFound, "user %d not found", req.GetId())
	}
	return entity.ToProto(), nil
}

func (s *Service) DeleteUser(ctx context.Context, req *commonv1.Identity) (*empty.Empty, error) {
	if err := s.store.DeleteUser(ctx, req.GetId()); err != nil {
		return nil, status.Errorf(codes.Internal, "delete user %d: %v", req.GetId(), err)
	}
	return &empty.Empty{}, nil
}
