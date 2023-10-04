package grpc

import (
	"github.com/kolesnikovm/messenger/proto"
	"google.golang.org/grpc"
)

type ServerBuilder struct {
	MessengerServer proto.MessengerServer
	Interceptor     grpc.StreamServerInterceptor
}

func (s *ServerBuilder) Build() *grpc.Server {
	srv := grpc.NewServer(
		grpc.StreamInterceptor(s.Interceptor),
	)
	proto.RegisterMessengerServer(srv, s.MessengerServer)

	return srv
}
