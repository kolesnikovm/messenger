package di

import (
	"github.com/kolesnikovm/messenger/configs"
	"github.com/kolesnikovm/messenger/notifier"
	"github.com/kolesnikovm/messenger/notifier/kafka"
)

func ProvideNotifier(conf configs.ServerConfig) (notifier.MessageSender, error) {
	return kafka.New(conf.KafkaConfig)
}
