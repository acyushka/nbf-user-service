package grpc_server

import (
	"context"
	"fmt"
	"nbf-user/internal/app"
	"nbf-user/internal/config"
	"nbf-user/internal/repository/postgres"
	"net"

	"github.com/hesoyamTM/nbf-auth/pkg/logger"
	userv1 "github.com/hesoyamTM/nbf-protos/gen/go/user"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GrpcServer struct {
	server *grpc.Server
	host   string
	port   int
}

func NewGrpcServer(ctx context.Context, db *sqlx.DB, cfg *config.Config) (*GrpcServer, error) {
	const op = "grpc.NewGrpcServer"

	// Init bd
	userRepo := postgres.NewUserRepo(db)

	// Init UserService
	//там логика всего микросервиса
	userService := app.NewUserService(userRepo)

	// Init UserServer
	// Там обертка под интерфейс юзер сервиса
	userServer := NewUserServer(userService)

	// Init interceptor
	logInterceptor, err := logger.NewLoggingInterceptor(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Create grpc server
	server := grpc.NewServer(grpc.UnaryInterceptor(logInterceptor))

	//Register
	userv1.RegisterUserServer(server, userServer)

	//Reflection
	reflection.Register(server)

	return &GrpcServer{
		server: server,
		port:   cfg.Port,
		host:   cfg.Host,
	}, nil
}

func (s *GrpcServer) MustStart(ctx context.Context) {
	const op = "grpc.MustStart"

	log, err := logger.LoggerFromCtx(ctx)
	if err != nil {
		panic(fmt.Errorf("%s: %w", op, err))
	}

	log.Info("grpc server is starting")

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.host, s.port))
	if err != nil {
		panic(fmt.Errorf("%s:%w", op, err))
	}

	log.Info("Server is started")

	if err := s.server.Serve(lis); err != nil {
		panic(fmt.Errorf("%s:%w", op, err))
	}

}

func (s *GrpcServer) MustStop(ctx context.Context) {
	const op = "grpc.MustStop"

	log, err := logger.LoggerFromCtx(ctx)
	if err != nil {
		panic(fmt.Errorf("%s: %w", op, err))
	}

	log.Info("grpc server is stopping")

	s.server.GracefulStop()
}
