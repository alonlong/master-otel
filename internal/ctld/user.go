package ctld

import (
	"context"

	commonv1 "master-otel/internal/proto/common/v1"

	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Service) CreateUser(ctx context.Context, req *commonv1.User) (*commonv1.User, error) {
	return s.storedClient.CreateUser(ctx, req)
}

func (s *Service) GetUser(ctx context.Context, req *commonv1.Identity) (*commonv1.User, error) {
	return s.storedClient.GetUser(ctx, req)
}

func (s *Service) DeleteUser(ctx context.Context, req *commonv1.Identity) (*empty.Empty, error) {
	return s.storedClient.DeleteUser(ctx, req)
}
