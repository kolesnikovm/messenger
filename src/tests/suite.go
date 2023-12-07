package tests

import (
	"context"
	"net"
	"testing"

	"github.com/kolesnikovm/messenger/archiver"
	notifier "github.com/kolesnikovm/messenger/notifier/mocks"
	"github.com/kolesnikovm/messenger/proto"
	store "github.com/kolesnikovm/messenger/store/mocks"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

type Suite struct {
	grpcServer             *grpc.Server
	messengerServiceClient proto.MessengerClient
	conn                   *grpc.ClientConn
	messageSender          *notifier.MockMessageSender
	messageStore           *store.MockMessages
	archiver               archiver.Archiver
	t                      *testing.T
}

func newConnection(t *testing.T, grpcServer *grpc.Server) (*grpc.ClientConn, error) {
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

	return conn, nil
}

func newClient(conn *grpc.ClientConn) proto.MessengerClient {
	return proto.NewMessengerClient(conn)
}
