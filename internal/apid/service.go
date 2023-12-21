package apid

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	storedv1 "master-otel/internal/proto/stored/v1"
	"master-otel/internal/utils"
	"master-otel/pkg/log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Service struct {
	e *echo.Echo

	storedAddr   string
	storedConn   *grpc.ClientConn
	storedClient storedv1.StoredServiceClient
}

func NewService(httpAddr string, ctldAddr string, service string) *Service {
	s := &Service{
		e:          echo.New(),
		storedAddr: ctldAddr,
	}
	s.e.Server.Addr = httpAddr
	s.e.Logger.SetOutput(io.Discard)
	s.e.HideBanner = true
	s.e.HidePort = true
	s.initRoutes(service)
	return s
}

func (s *Service) initRoutes(service string) {
	s.e.Use(middleware.Recover())
	s.e.Use(utils.TraceMiddleware(service))
	s.e.Use(utils.LoggerMiddleware())

	s.e.POST("/user", s.createUser)
	s.e.DELETE("/user/:id", s.deleteUser)
}

func (s *Service) Run(ctx context.Context) error {
	storedConn, err := utils.GrpcDial(ctx, s.storedAddr)
	if err != nil {
		return fmt.Errorf("dial ctld %s: %w", s.storedAddr, err)
	}
	s.storedConn = storedConn
	s.storedClient = storedv1.NewStoredServiceClient(storedConn)
	log.Info("connect to ctld service", zap.String("addr", s.storedAddr))

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
	if s.storedConn != nil {
		if err := s.storedConn.Close(); err != nil {
			log.Error("close ctld connection", zap.Error(err))
		}
	}
}
