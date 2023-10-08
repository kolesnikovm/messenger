package tests

import (
	"google.golang.org/grpc"
)

type Suite struct {
	grpcServer *grpc.Server
}

func newSuite(grpcServer *grpc.Server) *Suite {
	return &Suite{
		grpcServer: grpcServer,
	}
}
