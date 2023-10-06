//go:build wireinject
// +build wireinject

package grpc

import (
	"github.com/google/wire"
	"github.com/kolesnikovm/messenger/proto"
	"github.com/kolesnikovm/messenger/server/grpc/messenger"
	"github.com/kolesnikovm/messenger/usecase"
	"github.com/kolesnikovm/messenger/usecase/message"
)

func InitializeServerBuilder() ServerBuilder {
	wire.Build(
		message.New,
		wire.Bind(new(usecase.Message), new(*message.MessageUseCase)),
		messenger.NewHandler,
		wire.Bind(new(proto.MessengerServer), new(*messenger.Handler)),
		NewInterceptor,
		wire.Struct(new(ServerBuilder), "*"),
	)
	return ServerBuilder{}
}
