package grpc

import (
	"github.com/kolesnikovm/messenger/proto"
	"google.golang.org/grpc"
)

type ServerBuilder struct {
	MessengerServer   proto.MessengerServer
	StreamInterceptor grpc.StreamServerInterceptor
	UnaryInterceptor  grpc.UnaryServerInterceptor
}

func (s *ServerBuilder) Build() *grpc.Server {
	srv := grpc.NewServer(
		grpc.StreamInterceptor(s.StreamInterceptor),
		grpc.UnaryInterceptor(s.UnaryInterceptor),
	)
	proto.RegisterMessengerServer(srv, s.MessengerServer)

	return srv
}
