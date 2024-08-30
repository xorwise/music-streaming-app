package authgrpc

import (
	"context"
	"errors"

	sso "github.com/xorwise/music-streaming-service/gen"
	database "github.com/xorwise/music-streaming-service/internal/database"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	sso.UnimplementedAuthServer
	auth Auth
}

type Auth interface {
	Login(ctx context.Context, username string, password string) (token string, err error)
	Register(ctx context.Context, username string, password string) (userID int64, err error)
}

func Register(gRPCServer *grpc.Server, auth Auth) {
	sso.RegisterAuthServer(gRPCServer, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, req *sso.LoginRequest) (*sso.LoginResponse, error) {
	if req.GetUsername() == "" || req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "username and password are required")
	}

	token, err := s.auth.Login(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {

		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &sso.LoginResponse{Token: token}, nil

}

func (s *serverAPI) Register(ctx context.Context, req *sso.RegisterRequest) (*sso.RegisterResponse, error) {
	if req.GetUsername() == "" || req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "username and password are required")
	}
	userID, err := s.auth.Register(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		if errors.Is(err, database.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "failed to register")
	}
	return &sso.RegisterResponse{Id: userID}, nil
}
