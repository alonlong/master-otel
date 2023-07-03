package main

import (
	"context"
	"flag"
	"net"
	"os"
	"os/signal"
	"syscall"

	storedv1 "master-otel/internal/proto/stored/v1"
	"master-otel/internal/stored"
	"master-otel/internal/utils"
	"master-otel/pkg/db"
	"master-otel/pkg/log"

	"github.com/joho/godotenv"
	"github.com/uptrace/opentelemetry-go-extra/otelplay"
	"go.uber.org/zap"
)

var (
	grpcAddr = flag.String("grpc-addr", "localhost:9082", "The grpc address to bind")
)

func main() {
	flag.Parse()

	godotenv.Load(".env")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	otelShutdown := otelplay.ConfigureOpentelemetry(context.Background())
	defer otelShutdown()

	ln, err := net.Listen("tcp", *grpcAddr)
	if err != nil {
		log.Fatal("net listen", zap.String("grpc-addr", *grpcAddr), zap.Error(err))
	}

	dbConfig := &db.Config{
		Host:     os.Getenv("POSTGRES_HOST"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DB:       os.Getenv("POSTGRES_DB"),
		Port:     os.Getenv("POSTGRES_PORT"),
	}
	storedService, err := stored.NewService(&stored.Config{
		DB: dbConfig,
	})
	if err != nil {
		log.Fatal("new stored server", zap.Error(err))
	}
	storedService.Run(ctx)

	gs := utils.NewGrpcServer()
	storedv1.RegisterStoredServiceServer(gs, storedService)
	go func() {
		log.Info("start grpc server", zap.String("grpc-addr", *grpcAddr))
		if err := gs.Serve(ln); err != nil {
			log.Fatal("grpc serve", zap.Error(err))
		}
	}()

	<-ctx.Done()

	// stop the service
	storedService.Shutdown()
	// gracefully stop the grpc server
	gs.GracefulStop()
}
