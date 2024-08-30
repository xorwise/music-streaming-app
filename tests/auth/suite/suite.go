package suite

import (
	"context"
	"net"
	"strconv"
	"testing"

	ssov1 "github.com/xorwise/music-streaming-service/gen"
	config "github.com/xorwise/music-streaming-service/internal/config/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient ssov1.AuthClient
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.New()

	ctx, cancel := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancel()
	})

	grpcAddress := net.JoinHostPort("127.0.0.1", strconv.Itoa(cfg.GRPC.Port))
	cc, err := grpc.NewClient(
		grpcAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	authClient := ssov1.NewAuthClient(cc)

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: authClient,
	}

}
