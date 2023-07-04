package main

import (
	"context"
	"flag"
	"net"
	"os/signal"
	"syscall"

	"master-otel/internal/ctld"
	ctldv1 "master-otel/internal/proto/ctld/v1"
	"master-otel/internal/utils"
	"master-otel/pkg/log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	grpcAddr   = flag.String("grpc-addr", "localhost:9083", "The grpc address to bind")
	storedAddr = flag.String("stored-addr", "localhost:9082", "The stored server's address")
)

func main() {
	flag.Parse()

	logger := log.Init(&log.Config{Filename: "logs/ctld.log", MinLevel: zapcore.InfoLevel, Stdout: true})
	defer logger.Sync()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	grpcListener, err := net.Listen("tcp", *grpcAddr)
	if err != nil {
		log.Fatal("net listen", zap.String("grpc-addr", *grpcAddr), zap.Error(err))
	}
	ctldService := ctld.NewService(*storedAddr)
	if err := ctldService.Run(ctx); err != nil {
		log.Fatal("run ctld server", zap.Error(err))
	}

	gs := utils.NewGrpcServer("ctld")
	ctldv1.RegisterCtldServiceServer(gs, ctldService)
	go func() {
		log.Info("start grpc server", zap.String("grpc-addr", *grpcAddr))
		if err := gs.Serve(grpcListener); err != nil {
			log.Fatal("grpc serve", zap.Error(err))
		}
	}()

	<-ctx.Done()
	// stop the service
	ctldService.Shutdown()
	// gracefully stop the grpc server
	gs.GracefulStop()
}
