package tests

import (
	"context"
	"net"
	"testing"

	"github.com/kolesnikovm/messenger/notifier/mocks"
	"github.com/kolesnikovm/messenger/proto"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

type Suite struct {
	grpcServer             *grpc.Server
	messengerServiceClient proto.MessengerClient
	conn                   *grpc.ClientConn
	messageSender          *mocks.MockMessageSender
	t                      *testing.T
}

func newSuite(t *testing.T, grpcServer *grpc.Server, messageSender *mocks.MockMessageSender) (*Suite, error) {
	const bufSize = 1024 * 1024
	lis := bufconn.Listen(bufSize)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			require.NoErrorf(t, err, "server exited with error")
		}
	}()

	bufDialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	conn, err := grpc.Dial("bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		grpcServer.Stop()
		return nil, err
	}

	messengerServiceClient := proto.NewMessengerClient(conn)

	return &Suite{
		grpcServer:             grpcServer,
		messengerServiceClient: messengerServiceClient,
		conn:                   conn,
		messageSender:          messageSender,
		t:                      t,
	}, nil
}

func (s *Suite) Stop() {
	s.grpcServer.Stop()

	err := s.conn.Close()
	require.NoErrorf(s.t, err, "failed to close grpc.ClientConn")
}
