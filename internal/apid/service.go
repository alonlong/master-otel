package apid

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	ctldv1 "master-otel/internal/proto/ctld/v1"
	"master-otel/internal/utils"
	"master-otel/pkg/log"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Service struct {
	e *echo.Echo

	ctldAddr   string
	ctldConn   *grpc.ClientConn
	ctldClient ctldv1.CtldServiceClient
}

func NewService(httpAddr string, ctldAddr string) *Service {
	s := &Service{
		e:        echo.New(),
		ctldAddr: ctldAddr,
	}
	s.e.Server.Addr = httpAddr
	s.e.Logger.SetOutput(io.Discard)
	s.e.HideBanner = true
	s.e.HidePort = true
	s.initRoutes()
	return s
}

func (s *Service) initRoutes() {
	s.e.POST("/user", s.createUser)
}

func (s *Service) Run(ctx context.Context) error {
	ctldConn, err := utils.GrpcDial(ctx, s.ctldAddr)
	if err != nil {
		return fmt.Errorf("dial ctld %s: %w", s.ctldAddr, err)
	}
	s.ctldConn = ctldConn
	s.ctldClient = ctldv1.NewCtldServiceClient(ctldConn)
	log.Info("connect to ctld service", zap.String("addr", s.ctldAddr))

	// start the apid server
	go func() {
		log.Info("start apid server", zap.String("addr", s.e.Server.Addr))
		if err := s.e.Start(s.e.Server.Addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("start apid server", zap.Error(err))
		}
	}()
	return nil
}

func (s *Service) Shutdown() {
	if err := s.e.Close(); err != nil {
		log.Error("shutdown apid server", zap.Error(err))
	}
	if s.ctldConn != nil {
		if err := s.ctldConn.Close(); err != nil {
			log.Error("close ctld connection", zap.Error(err))
		}
	}
}
