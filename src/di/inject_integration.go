package di

import (
	"github.com/google/wire"
	"github.com/kolesnikovm/messenger/configs"
	"github.com/kolesnikovm/messenger/notifier"
	"github.com/kolesnikovm/messenger/notifier/kafka"
)

func ProvideKafka(conf configs.ServerConfig) *kafka.KafkaMessageSender {
	return kafka.New(conf.KafkaConfig)
}

var NotifierSet = wire.NewSet(
	ProvideKafka,
	wire.Bind(new(notifier.MessageSender), new(*kafka.KafkaMessageSender)),
)
