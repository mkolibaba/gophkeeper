package grpc

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	pb "github.com/mkolibaba/gophkeeper/internal/common/grpc/proto/gen"
	"github.com/mkolibaba/gophkeeper/internal/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"sync"
)

type Server struct {
	s      *grpc.Server
	logger *zap.Logger
}

func NewServer(
	loginService server.LoginService,
	noteService server.NoteService,
	binaryService server.BinaryService,
	cardService server.CardService,
	dataValidator *validator.Validate,
	logger *zap.Logger,
) *Server {
	dataServiceServer := &dataServiceServer{
		loginService:  loginService,
		noteService:   noteService,
		binaryService: binaryService,
		cardService:   cardService,
		dataValidator: dataValidator,
		logger:        logger,
	}

	s := grpc.NewServer()
	pb.RegisterDataServiceServer(s, dataServiceServer)
	reflection.Register(s)

	return &Server{
		s:      s,
		logger: logger,
	}
}

func (s *Server) Start(ctx context.Context, addr string) error {
	s.logger.Info("running grpc server", zap.String("addr", addr))

	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	var wg sync.WaitGroup
	wg.Go(func() {
		<-ctx.Done()
		s.s.GracefulStop()
	})

	if err := s.s.Serve(listen); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	s.logger.Info("server stopped")

	wg.Wait()
	return nil
}
