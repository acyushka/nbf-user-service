package main

import (
	"context"
	"nbf-user/internal/config"
	grpc_server "nbf-user/internal/server/grpc"
	"os"
	"os/signal"
	"syscall"

	cfgtools "github.com/hesoyamTM/nbf-auth/pkg/config"
	"github.com/hesoyamTM/nbf-auth/pkg/logger"

	"github.com/jmoiron/sqlx"
	_ "github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	cfg := cfgtools.MustParseConfig[config.Config]()
	ctx, err := logger.SetupLogger(context.Background(), cfg.Env)
	if err != nil {
		panic(err)
	}

	log, err := logger.LoggerFromCtx(ctx)
	if err != nil {
		panic(err)
	}

	log.Debug("Logger is working")

	db, err := sqlx.Connect("postgres", cfg.Database.GetDatabaseURL())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	gRPCserver, err := grpc_server.NewGrpcServer(ctx, db, cfg)
	if err != nil {
		panic(err)
	}

	go gRPCserver.MustStart(ctx)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	gRPCserver.MustStop(ctx)

	log.Info("grpc server is gracefully stopped")
}
