package grpc

import (
	"github.com/kolesnikovm/messenger/proto"
	"google.golang.org/grpc"
)

type ServerBuilder struct {
	MessengerServer proto.MessengerServer
}

func (s *ServerBuilder) Build() *grpc.Server {
	srv := grpc.NewServer()
	proto.RegisterMessengerServer(srv, s.MessengerServer)

	return srv
}
