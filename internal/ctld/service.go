package ctld

import (
	"context"
	"fmt"

	ctldv1 "master-otel/internal/proto/ctld/v1"
	storedv1 "master-otel/internal/proto/stored/v1"
	"master-otel/internal/utils"
	"master-otel/pkg/log"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Service struct {
	ctldv1.UnimplementedCtldServiceServer

	storedAddr   string
	storedConn   *grpc.ClientConn
	storedClient storedv1.StoredServiceClient
}

func NewService(storedAddr string) *Service {
	return &Service{
		storedAddr: storedAddr,
	}
}

func (s *Service) Run(ctx context.Context) error {
	storedConn, err := utils.GrpcDial(ctx, s.storedAddr)
	if err != nil {
		return fmt.Errorf("dial stored %s: %w", s.storedAddr, err)
	}
	s.storedConn = storedConn
	s.storedClient = storedv1.NewStoredServiceClient(storedConn)
	log.Info("connect to stored service", zap.String("addr", s.storedAddr))
	return nil
}

func (s *Service) Shutdown() {
	if s.storedConn != nil {
		if err := s.storedConn.Close(); err != nil {
			log.Error("close stored connection", zap.Error(err))
		}
	}
}
