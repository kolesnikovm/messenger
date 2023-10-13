package di

import (
	"github.com/google/wire"
	"github.com/kolesnikovm/messenger/configs"
	"github.com/kolesnikovm/messenger/integration"
	"github.com/kolesnikovm/messenger/integration/kafka"
)

func ProvideKafka(conf configs.ServerConfig) *kafka.KafkaMessageSender {
	return kafka.New(conf.KafkaConfig)
}

var IntegrationSet = wire.NewSet(
	ProvideKafka,
	wire.Bind(new(integration.MessageSender), new(*kafka.KafkaMessageSender)),
)
