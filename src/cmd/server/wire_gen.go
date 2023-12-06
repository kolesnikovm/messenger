// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package server

import (
	"github.com/kolesnikovm/messenger/configs"
	"github.com/kolesnikovm/messenger/di"
	"github.com/kolesnikovm/messenger/server/grpc"
	"github.com/kolesnikovm/messenger/server/grpc/messenger"
	"github.com/kolesnikovm/messenger/usecase/message"
)

// Injectors from wire.go:

func InitializeApplication(conf configs.ServerConfig) (*application, func(), error) {
	messageSender, cleanup, err := di.ProvideNotifier(conf)
	if err != nil {
		return nil, nil, err
	}
	db, cleanup2, err := di.ProvideDB(conf)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	messages := di.ProvideMessages(db, conf)
	orderIDCacher, cleanup3, err := di.ProvideCache(conf)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	messageUseCase := &message.MessageUseCase{
		MessageSender: messageSender,
		MessageStore:  messages,
		OrderIDCache:  orderIDCacher,
	}
	handler := messenger.NewHandler(messageUseCase)
	serverBuilder := grpc.ServerBuilder{
		MessengerServer: handler,
	}
	server := di.ProvideServer(serverBuilder)
	archiver, cleanup4, err := di.ProvideArchiver(conf, messages)
	if err != nil {
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	serverApplication := &application{
		grpcServer: server,
		archiver:   archiver,
		notifier:   messageSender,
	}
	return serverApplication, func() {
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}
