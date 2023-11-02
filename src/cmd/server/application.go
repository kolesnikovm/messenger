package server

import (
	"github.com/kolesnikovm/messenger/archiver"
	"github.com/kolesnikovm/messenger/notifier"
	"google.golang.org/grpc"
)

type application struct {
	grpcServer *grpc.Server
	archiver   archiver.Archiver
	notifier   notifier.MessageSender
}
