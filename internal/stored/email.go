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

func (s *Service) CreateEmail(ctx context.Context, req *commonv1.Email) (*commonv1.Email, error) {
	log.WithCtx(ctx).Info("create email", zap.String("email", req.GetEmail()))
	entity := models.Email{
		Email: req.GetEmail(),
	}
	if err := s.store.CreateEmail(ctx, &entity); err != nil {
		return nil, status.Errorf(codes.Internal, "create email %s: %v", req.GetEmail(), err)
	}
	return entity.ToProto(), nil
}
