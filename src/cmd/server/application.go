package server

import (
	"github.com/kolesnikovm/messenger/archiver"
	"github.com/kolesnikovm/messenger/store"
	"google.golang.org/grpc"
)

type application struct {
	grpcServer *grpc.Server
	archiver   archiver.Archiver
	aggregator store.Aggregator
}
