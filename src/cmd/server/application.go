package server

import (
	"google.golang.org/grpc"
)

type application struct {
	grpcServer *grpc.Server
}

func newApplication(grpcServer *grpc.Server) *application {
	return &application{
		grpcServer: grpcServer,
	}
}
