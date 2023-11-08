package di

import (
	"github.com/google/wire"
	"github.com/kolesnikovm/messenger/usecase"
	"github.com/kolesnikovm/messenger/usecase/message"
)

var UsecaseSet = wire.NewSet(
	wire.Struct(new(message.MessageUseCase), "*"),
	wire.Bind(new(usecase.Message), new(*message.MessageUseCase)),
)
