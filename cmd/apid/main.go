package main

import (
	"context"
	"flag"
	"os/signal"
	"syscall"

	"master-otel/internal/apid"
	"master-otel/pkg/log"

	"go.uber.org/zap"
)

var (
	httpAddr = flag.String("http-addr", "localhost:9084", "The http server's address")
	ctldAddr = flag.String("ctld-addr", "localhost:9083", "The ctld server's address")
)

func main() {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	apidService := apid.NewService(*httpAddr, *ctldAddr)
	if err := apidService.Run(ctx); err != nil {
		log.Fatal("run apid server", zap.Error(err))
	}
	<-ctx.Done()

	// stop the service
	apidService.Shutdown()
}
