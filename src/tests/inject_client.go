package tests

import (
	"testing"

	"github.com/google/wire"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func ProvideConnection(t *testing.T, grpcServer *grpc.Server) (*grpc.ClientConn, func(), error) {
	conn, err := newConnection(t, grpcServer)

	cleanup := func() {
		grpcServer.Stop()

		err := conn.Close()
		require.NoErrorf(t, err, "failed to close grpc.ClientConn")
	}

	return conn, cleanup, err
}

var ClientSet = wire.NewSet(
	ProvideConnection,
	newClient,
)
