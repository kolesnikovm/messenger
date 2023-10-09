package tests

import (
	"context"
	"net"

	"github.com/kolesnikovm/messenger/proto"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

type Suite struct {
	grpcServer             *grpc.Server
	messengerServiceClient proto.MessengerClient
	conn                   *grpc.ClientConn
}

func newSuite(grpcServer *grpc.Server) (*Suite, error) {
	const bufSize = 1024 * 1024
	lis := bufconn.Listen(bufSize)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal().Err(err).Msg("server exited with error")
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
	}, nil
}

func (s *Suite) Stop() {
	if err := s.conn.Close(); err != nil {
		log.Fatal().Err(err).Msg("failed to close grpc.ClientConn")
	}
	s.grpcServer.Stop()
}
