package di

import (
	"github.com/google/wire"
	"github.com/kolesnikovm/messenger/proto"
	grpcServer "github.com/kolesnikovm/messenger/server/grpc"
	"github.com/kolesnikovm/messenger/server/grpc/messenger"
	"google.golang.org/grpc"
)

func ProvideServer(builder grpcServer.ServerBuilder) *grpc.Server {
	return builder.Build()
}

var ServerSet = wire.NewSet(
	messenger.NewHandler,
	wire.Bind(new(proto.MessengerServer), new(*messenger.Handler)),
	grpcServer.NewInterceptor,
	wire.Struct(new(grpcServer.ServerBuilder), "*"),
	ProvideServer,
)
