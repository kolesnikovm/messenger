package tests

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/kolesnikovm/messenger/proto"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	suite := InitializeSuite()
	go func() {
		if err := suite.grpcServer.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestSendMessage(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := proto.NewMessengerClient(conn)
	stream, err := client.SendMessage(ctx)
	if err != nil {
		t.Fatalf("Failed to create stream: %v", err)
	}

	message := &proto.Message{Text: "test"}
	if err := stream.Send(message); err != nil {
		t.Fatalf("%v.Send(%v) = %v", stream, message, err)
	}

	reply, err := stream.CloseAndRecv()
	if err != nil {
		t.Fatalf("%v.CloseAndRecv() got error %v", stream, err)
	}

	require.Equal(t, 0, int(reply.GetErrorCount()))
}
