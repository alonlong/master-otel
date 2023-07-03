package stored

import (
	"context"

	"master-otel/internal/entity/models"
	storedv1 "master-otel/internal/proto/stored/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) CreateUser(ctx context.Context, req *storedv1.CreateUserRequest) (*storedv1.CreateUserResponse, error) {
	entity := models.User{
		Email:    req.GetUser().GetEmail(),
		Username: req.GetUser().GetUsername(),
	}
	if err := s.store.CreateUser(&entity); err != nil {
		return nil, status.Errorf(codes.Internal, "create user %s: %v", req.GetUser().GetEmail(), err)
	}
	return &storedv1.CreateUserResponse{
		Id: entity.ID,
	}, nil
}

func (s *Service) GetUser(ctx context.Context, req *storedv1.GetUserRequest) (*storedv1.GetUserResponse, error) {
	entity, err := s.store.GetUser(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "query user %d: %v", req.GetId(), err)
	}
	if entity == nil {
		return nil, status.Errorf(codes.NotFound, "user %d not found", req.GetId())
	}
	return &storedv1.GetUserResponse{
		User: entity.ToProto(),
	}, nil
}

func (s *Service) DeleteUser(ctx context.Context, req *storedv1.DeleteUserRequest) (*storedv1.DeleteUserResponse, error) {
	if err := s.store.DeleteUser(req.GetId()); err != nil {
		return nil, status.Errorf(codes.Internal, "delete user %d: %v", req.GetId(), err)
	}
	return &storedv1.DeleteUserResponse{}, nil
}
