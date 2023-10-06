package server

import "github.com/kolesnikovm/messenger/server/grpc"

type application struct {
	grpcServerBuilder *grpc.ServerBuilder
}

func newApplication(grpcServerBuilder *grpc.ServerBuilder) *application {
	return &application{
		grpcServerBuilder: grpcServerBuilder,
	}
}
