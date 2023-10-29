package di

import (
	"github.com/kolesnikovm/messenger/configs"
	"github.com/kolesnikovm/messenger/notifier"
	"github.com/kolesnikovm/messenger/notifier/kafka"
	"github.com/kolesnikovm/messenger/store"
)

func ProvideNotifier(conf configs.ServerConfig, messageStore store.Messages) (notifier.MessageSender, func(), error) {
	messageSender, err := kafka.New(conf.KafkaConfig, messageStore)

	cleanup := func() {
		messageSender.Close()
	}

	return messageSender, cleanup, err
}
