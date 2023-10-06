package di

import (
	"github.com/google/wire"
	"github.com/kolesnikovm/messenger/proto"
	"github.com/kolesnikovm/messenger/server/grpc"
	"github.com/kolesnikovm/messenger/server/grpc/messenger"
)

var ServerSet = wire.NewSet(
	messenger.NewHandler,
	wire.Bind(new(proto.MessengerServer), new(*messenger.Handler)),
	grpc.NewInterceptor,
	wire.Struct(new(grpc.ServerBuilder), "*"),
)
